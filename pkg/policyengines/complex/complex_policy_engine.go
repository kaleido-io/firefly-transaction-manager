// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package complex

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/ffresty"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmmsgs"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/hyperledger/firefly-transaction-manager/pkg/policyengine"
)

type PolicyEngineFactory struct{}

func (f *PolicyEngineFactory) Name() string {
	return "complex"
}

// simplePolicyEngine is a base policy engine forming an example for extension:
// - It uses a public gas estimation
// - It submits the transaction once
// - It logs errors transactions breach certain configured thresholds of staleness
func (f *PolicyEngineFactory) NewPolicyEngine(ctx context.Context, conf config.Section) (pe policyengine.PolicyEngine, err error) {
	gasOracleConfig := conf.SubSection(GasOracleConfig)
	p := &complexPolicyEngine{
		resubmitInterval: conf.GetDuration(ResubmitInterval),
		fixedGasPrice:    fftypes.JSONAnyPtr(conf.GetString(FixedGasPrice)),

		gasOracleMethod:        gasOracleConfig.GetString(GasOracleMethod),
		gasOracleQueryInterval: gasOracleConfig.GetDuration(GasOracleQueryInterval),
		gasOracleMode:          gasOracleConfig.GetString(GasOracleMode),
	}
	log.L(ctx).Info("Starting complex policy engine")
	switch p.gasOracleMode {
	case GasOracleModeConnector:
		// No initialization required
	case GasOracleModeRESTAPI:
		p.gasOracleClient = ffresty.New(ctx, gasOracleConfig)
		templateString := gasOracleConfig.GetString(GasOracleTemplate)
		if templateString == "" {
			return nil, i18n.NewError(ctx, tmmsgs.MsgMissingGOTemplate)
		}
		p.gasOracleTemplate, err = template.New("").Funcs(sprig.FuncMap()).Parse(templateString)
		if err != nil {
			return nil, i18n.NewError(ctx, tmmsgs.MsgBadGOTemplate, err)
		}
	default:
		if p.fixedGasPrice.IsNil() {
			return nil, i18n.NewError(ctx, tmmsgs.MsgNoGasConfigSetForPolicyEngine)
		}
	}
	return p, nil
}

type complexPolicyEngine struct {
	fixedGasPrice    *fftypes.JSONAny
	resubmitInterval time.Duration

	gasOracleMode          string
	gasOracleClient        *resty.Client
	gasOracleMethod        string
	gasOracleTemplate      *template.Template
	gasOracleQueryInterval time.Duration
	gasOracleQueryValue    *fftypes.JSONAny
	gasOracleLastQueryTime *fftypes.FFTime
}

type complexPolicyInfo struct {
	LastWarnTime *fftypes.FFTime `json:"lastWarnTime"`
}

// withPolicyInfo is a convenience helper to run some logic that accesses/updates our policy section
func (p *complexPolicyEngine) withPolicyInfo(ctx context.Context, mtx *apitypes.ManagedTX, fn func(info *complexPolicyInfo) (update policyengine.UpdateType, result policyengine.PolicyExecutionResult, reason ffcapi.ErrorReason, err error)) (update policyengine.UpdateType, result policyengine.PolicyExecutionResult, reason ffcapi.ErrorReason, err error) {
	var info complexPolicyInfo
	infoBytes := []byte(mtx.PolicyInfo.String())
	if len(infoBytes) > 0 {
		err := json.Unmarshal(infoBytes, &info)
		if err != nil {
			log.L(ctx).Warnf("Failed to parse existing info `%s`: %s", infoBytes, err)
		}
	}
	update, result, reason, err = fn(&info)
	if update != policyengine.UpdateNo {
		infoBytes, _ = json.Marshal(&info)
		mtx.PolicyInfo = fftypes.JSONAnyPtrBytes(infoBytes)
	}
	return update, result, reason, err
}

func (p *complexPolicyEngine) submitTX(ctx context.Context, cAPI ffcapi.API, mtx *apitypes.ManagedTX) (reason ffcapi.ErrorReason, err error) {
	sendTX := &ffcapi.TransactionSendRequest{
		TransactionHeaders: mtx.TransactionHeaders,
		GasPrice:           mtx.GasPrice,
		TransactionData:    mtx.TransactionData,
	}
	sendTX.TransactionHeaders.Nonce = (*fftypes.FFBigInt)(mtx.Nonce.Int())
	sendTX.TransactionHeaders.Gas = (*fftypes.FFBigInt)(mtx.Gas.Int())
	log.L(ctx).Debugf("Complex engine sending transaction %s at nonce %s / %d (lastSubmit=%s)", mtx.ID, mtx.TransactionHeaders.From, mtx.Nonce.Int64(), mtx.LastSubmit)
	res, reason, err := cAPI.TransactionSend(ctx, sendTX)
	if err == nil {
		mtx.TransactionHash = res.TransactionHash
		mtx.LastSubmit = fftypes.Now()
	} else {
		// We have some simple rules for handling reasons from the connector, which could be enhanced by extending the connector.
		switch reason {
		case ffcapi.ErrorKnownTransaction:
			// If we already have a transaction hash, this is fine - we just return as if we submitted it
			if mtx.TransactionHash != "" {
				log.L(ctx).Debugf("Complex engine transaction %s at nonce %s / %d known with hash: %s (%s)", mtx.ID, mtx.TransactionHeaders.From, mtx.Nonce.Int64(), mtx.TransactionHash, err)
				return "", nil
			}
			// Note: to cover the edge case where we had a timeout or other failure during the initial TransactionSend,
			//       a policy engine implementation would need to be able to re-calculate the hash that we would expect for the transaction.
			//       This would require a new FFCAPI API to calculate that hash, which requires the connector to perform the signing
			//       without submission to the node. For example using `eth_signTransaction` for EVM JSON/RPC.
			return reason, err
		case ffcapi.ErrorReasonNonceTooLow:
			return reason, err
		default:
			return reason, err
		}
	}
	log.L(ctx).Infof("Complex engine transaction %s at nonce %s / %d submitted. Hash: %s", mtx.ID, mtx.TransactionHeaders.From, mtx.Nonce.Int64(), mtx.TransactionHash)
	return "", nil
}

func (p *complexPolicyEngine) Execute(ctx context.Context, cAPI ffcapi.API, mtx *apitypes.ManagedTX) (update policyengine.UpdateType, result policyengine.PolicyExecutionResult, reason ffcapi.ErrorReason, err error) {

	// Simply policy engine allows deletion of the transaction without additional checks ( ensuring the TX has not been submitted / gap filling the nonce etc. )
	if mtx.DeleteRequested != nil {
		return policyengine.UpdateDelete, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, "", nil
	}

	// Simple policy engine only submits once.
	if mtx.FirstSubmit == nil {
		// Only calculate gas price here in the simple policy engine
		mtx.GasPrice, err = p.getGasPrice(ctx, cAPI)
		if err != nil {
			return policyengine.UpdateNo, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, "", err
		}
		// Submit the first time
		if reason, err := p.submitTX(ctx, cAPI, mtx); err != nil {
			if reason == ffcapi.ErrorReasonNonceTooLow {
				log.L(ctx).Infof("Complex engine nonce %s / %d for transaction %s is too low. Incrementing nonce and re-submitting.", mtx.TransactionHeaders.From, mtx.Nonce.Int64(), mtx.ID)
				mtx.Nonce = fftypes.NewFFBigInt(mtx.Nonce.Int64() + 1)

				resubReason, resubErr := p.submitTX(ctx, cAPI, mtx)
				if resubErr != nil {
					if resubReason == ffcapi.ErrorReasonNonceTooLow {
						log.L(ctx).Infof("Failed on trying to resubmit with bigger nonce. Give up")
						// We tried once more - give up for now and return the best hint we can to the caller
						mtx.Nonce = fftypes.NewFFBigInt(mtx.Nonce.Int64() + 1)
						return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceHint(mtx.Nonce.Int64())}, resubReason, resubErr
					}
				} else {
					// We managed to submit it. We mustn't return an error but we can at least give the transaction manager
					// a hint that their current nonce is incorrect for their next transaction
					log.L(ctx).Infof("Complex engine succeeded in re-submitting transaction %s after increasing the nonce", mtx.ID)
					mtx.FirstSubmit = mtx.LastSubmit
					return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceHint(mtx.Nonce.Int64())}, "", nil
				}
			}
			return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, reason, err
		}
		mtx.FirstSubmit = mtx.LastSubmit
		return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, "", nil

	} else if mtx.Receipt == nil {

		// A more sophisticated policy engine would look at the reason for the lack of a receipt, and consider taking progressive
		// action such as increasing the gas cost slowly over time. This simple example shows how the policy engine
		// can use the FireFly core operation as a store for its historical state/decisions (in this case the last time we warned).
		return p.withPolicyInfo(ctx, mtx, func(info *complexPolicyInfo) (update policyengine.UpdateType, result policyengine.PolicyExecutionResult, reason ffcapi.ErrorReason, err error) {
			lastWarnTime := info.LastWarnTime
			if lastWarnTime == nil {
				lastWarnTime = mtx.FirstSubmit
			}
			now := fftypes.Now()
			if now.Time().Sub(*lastWarnTime.Time()) > p.resubmitInterval {
				secsSinceSubmit := float64(now.Time().Sub(*mtx.FirstSubmit.Time())) / float64(time.Second)
				log.L(ctx).Infof("Complex engine transaction %s at nonce %s / %d has not been mined after %.2fs", mtx.ID, mtx.TransactionHeaders.From, mtx.Nonce.Int64(), secsSinceSubmit)
				info.LastWarnTime = now
				// We do a resubmit at this point - as it might no longer be in the TX pool
				if reason, err := p.submitTX(ctx, cAPI, mtx); err != nil {
					if reason == ffcapi.ErrorReasonNonceTooLow {
						// We seem to be out of step with the nonce for this signing address. Let's see if we can get the transaction submitted
						// and then return a hint to the nonce manager so it knows where the chain has got to
						log.L(ctx).Infof("Complex engine nonce %s / %d for transaction %s is too low. Incrementing nonce and re-submitting.", mtx.TransactionHeaders.From, mtx.Nonce.Int64(), mtx.ID)
						mtx.Nonce = fftypes.NewFFBigInt(mtx.Nonce.Int64() + 1)

						resubReason, resubErr := p.submitTX(ctx, cAPI, mtx)
						if resubErr != nil {
							if resubReason == ffcapi.ErrorReasonNonceTooLow {
								log.L(ctx).Infof("Failed on trying to resubmit with bigger nonce. Give up")
								// We tried once more - give up for now and return the best hint we can to the caller
								mtx.Nonce = fftypes.NewFFBigInt(mtx.Nonce.Int64() + 1)
								return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceHint(mtx.Nonce.Int64())}, resubReason, resubErr
							}
						} else {
							log.L(ctx).Infof("Complex engine succeeded in re-submitting transaction %s after increasing the nonce", mtx.ID)
							return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceHint(mtx.Nonce.Int64())}, "", nil
						}
					}
					if reason != ffcapi.ErrorKnownTransaction {
						return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, reason, err
					}
				}
				return policyengine.UpdateYes, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, "", nil
			}
			return policyengine.UpdateNo, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, "", nil
		})

	}
	// No action in the case we have a receipt
	return policyengine.UpdateNo, policyengine.PolicyExecutionResult{Hint: policyengine.NonceOK}, "", nil
}

// getGasPrice either uses a fixed gas price, or invokes a gas station API
func (p *complexPolicyEngine) getGasPrice(ctx context.Context, cAPI ffcapi.API) (gasPrice *fftypes.JSONAny, err error) {
	if p.gasOracleQueryValue != nil && p.gasOracleLastQueryTime != nil &&
		time.Since(*p.gasOracleLastQueryTime.Time()) < p.gasOracleQueryInterval {
		return p.gasOracleQueryValue, nil
	}
	switch p.gasOracleMode {
	case GasOracleModeRESTAPI:
		// Make a REST call against an endpoint, and extract a value/structure to pass to the connector
		gasPrice, err := p.getGasPriceAPI(ctx)
		if err != nil {
			return nil, err
		}
		p.gasOracleQueryValue = gasPrice
		p.gasOracleLastQueryTime = fftypes.Now()
		return p.gasOracleQueryValue, nil
	case GasOracleModeConnector:
		// Call the connector
		res, _, err := cAPI.GasPriceEstimate(ctx, &ffcapi.GasPriceEstimateRequest{})
		if err != nil {
			return nil, err
		}
		p.gasOracleQueryValue = res.GasPrice
		p.gasOracleLastQueryTime = fftypes.Now()
		return p.gasOracleQueryValue, nil
	default:
		// Disabled - just a fixed value - note that the fixed value can be any JSON structure,
		// as interpreted by the connector. For example EVMConnect support a simple value, or a
		// post EIP-1559 structure.
		return p.fixedGasPrice, nil
	}
}

func (p *complexPolicyEngine) getGasPriceAPI(ctx context.Context) (gasPrice *fftypes.JSONAny, err error) {
	var jsonResponse map[string]interface{}
	res, err := p.gasOracleClient.R().
		SetResult(&jsonResponse).
		Execute(p.gasOracleMethod, "")
	if err != nil {
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgErrorQueryingGasOracleAPI, -1, err.Error())
	}
	if res.IsError() {
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgErrorQueryingGasOracleAPI, res.StatusCode(), res.RawResponse)
	}
	buff := new(bytes.Buffer)
	err = p.gasOracleTemplate.Execute(buff, jsonResponse)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgGasOracleResultError)
	}
	return fftypes.JSONAnyPtr(buff.String()), nil
}
