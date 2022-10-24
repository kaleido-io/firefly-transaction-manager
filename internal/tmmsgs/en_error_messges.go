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

package tmmsgs

import (
	"net/http"

	"github.com/hyperledger/firefly-common/pkg/i18n"
	"golang.org/x/text/language"
)

var ffe = func(key, translation string, statusHint ...int) i18n.ErrorMessageKey {
	return i18n.FFE(language.AmericanEnglish, key, translation, statusHint...)
}

//revive:disable
var (
	MsgInvalidOutputType             = ffe("FF21010", "Invalid output type: %s")
	MsgConnectorError                = ffe("FF21012", "Connector failed request. requestId=%s reason=%s error: %s")
	MsgConnectorInvalidContentType   = ffe("FF21013", "Connector failed request. requestId=%s invalid response content type: %s")
	MsgInvalidConfirmationRequest    = ffe("FF21016", "Invalid confirmation request %+v")
	MsgCoreError                     = ffe("FF21017", "Error from core status=%d: %s")
	MsgConfigParamNotSet             = ffe("FF21018", "Configuration parameter '%s' must be set")
	MsgPolicyEngineNotRegistered     = ffe("FF21019", "No policy engine registered with name '%s'")
	MsgNoGasConfigSetForPolicyEngine = ffe("FF21020", "A fixed gas price must be set when not using a gas oracle")
	MsgErrorQueryingGasOracleAPI     = ffe("FF21021", "Error from gas station API [%d]: %s")
	MsgInvalidRequestErr             = ffe("FF21022", "Invalid '%s' request: %s", http.StatusBadRequest)
	MsgUnsupportedRequestType        = ffe("FF21023", "Unsupported request type: %s", http.StatusBadRequest)
	MsgMissingGOTemplate             = ffe("FF21024", "Missing template for processing response from Gas Oracle REST API")
	MsgBadGOTemplate                 = ffe("FF21025", "Invalid Go template: %s")
	MsgGasOracleResultError          = ffe("FF21026", "Error processing result from gas station API via template")
	MsgStreamStateError              = ffe("FF21027", "Event stream is in %s state", http.StatusConflict)
	MsgMissingName                   = ffe("FF21028", "Name is required", http.StatusBadRequest)
	MsgInvalidStreamType             = ffe("FF21029", "Invalid event stream type '%s'", http.StatusBadRequest)
	MsgMissingWebhookURL             = ffe("FF21030", "'url' is required for webhook configuration", http.StatusBadRequest)
	MsgStopFailedUpdatingESConfig    = ffe("FF21031", "Failed to stop event stream to apply updated configuration: %s")
	MsgStartFailedUpdatingESConfig   = ffe("FF21032", "Failed to restart event stream while applying updated configuration: %s")
	MsgBlockWebhookAddress           = ffe("FF21033", "Cannot send Webhook POST to address '%s' for host '%s'")
	MsgInvalidDistributionMode       = ffe("FF21034", "Invalid distribution mode for WebSocket: %s", http.StatusBadRequest)
	MsgWebhookFailedStatus           = ffe("FF21035", "Webhook request failed with status %d")
	MsgWSErrorFromClient             = ffe("FF21036", "Error received from WebSocket client: %s")
	MsgWebSocketClosed               = ffe("FF21037", "WebSocket '%s' closed")
	MsgWebSocketInterruptedSend      = ffe("FF21038", "Interrupted waiting for WebSocket connection to send event")
	MsgWebSocketInterruptedReceive   = ffe("FF21039", "Interrupted waiting for WebSocket acknowledgment")
	MsgBadListenerOptions            = ffe("FF21040", "Invalid listener options: %s", http.StatusBadRequest)
	MsgInvalidHost                   = ffe("FF21041", "Cannot send Webhook POST to host '%s': %s")
	MsgWebhookErr                    = ffe("FF21042", "Webhook request failed: %s")
	MsgUnknownPersistence            = ffe("FF21043", "Unknown persistence type '%s'")
	MsgInvalidLimit                  = ffe("FF21044", "Invalid limit string '%s': %s")
	MsgStreamNotFound                = ffe("FF21045", "Event stream '%v' not found", http.StatusNotFound)
	MsgListenerNotFound              = ffe("FF21046", "Event listener '%v' not found", http.StatusNotFound)
	MsgDuplicateStreamName           = ffe("FF21047", "Duplicate event stream name '%s' used by stream '%s'", http.StatusConflict)
	MsgMissingID                     = ffe("FF21048", "ID is required", http.StatusBadRequest)
	MsgPersistenceInitFail           = ffe("FF21049", "Failed to initialize '%s' persistence: %s")
	MsgLevelDBPathMissing            = ffe("FF21050", "Path must be supplied for LevelDB persistence")
	MsgFilterUpdateNotAllowed        = ffe("FF21051", "Event filters cannot be updated after a listener is created. Previous signature: '%s'. New signature: '%s'")
	MsgResetStreamNotFound           = ffe("FF21052", "Attempt to reset listener '%s', which is not currently registered on stream '%s'", http.StatusNotFound)
	MsgPersistenceMarshalFailed      = ffe("FF21053", "JSON serialization failed while writing to persistence")
	MsgPersistenceUnmarshalFailed    = ffe("FF21054", "JSON parsing failed while reading from persistence")
	MsgPersistenceReadFailed         = ffe("FF21055", "Failed to read key '%s' from persistence")
	MsgPersistenceWriteFailed        = ffe("FF21056", "Failed to write key from persistence")
	MsgPersistenceDeleteFailed       = ffe("FF21057", "Failed to delete from persistence")
	MsgPersistenceInitFailed         = ffe("FF21058", "Failed to initialize persistence at path '%s'")
	MsgPersistenceTXIncomplete       = ffe("FF21059", "Transaction is missing indexed fields")
	MsgNotStarted                    = ffe("FF21060", "Connector has not fully started yet", http.StatusServiceUnavailable)
	MsgPersistenceQueryFailed        = ffe("FF21061", "Persistence query failed")
	MsgPaginationErrTxNotFound       = ffe("FF21062", "The ID specified in the 'after' option (for pagination) must match an existing transaction: '%s'", http.StatusNotFound)
	MsgTXConflictSignerPending       = ffe("FF21063", "Only one of 'signer' and 'pending' can be supplied when querying transactions", http.StatusBadRequest)
	MsgInvalidSortDirection          = ffe("FF21064", "Sort direction must be 'asc'/'ascending' or 'desc'/'descending': '%s'", http.StatusBadRequest)
	MsgDuplicateID                   = ffe("FF21065", "ID '%s' is not unique", http.StatusConflict)
	MsgTransactionFailed             = ffe("FF21066", "Transaction execution failed")
	MsgTransactionNotFound           = ffe("FF21067", "Transaction '%s' not found", http.StatusNotFound)
	MsgPolicyEngineRequestTimeout    = ffe("FF21068", "The policy engine did not acknowledge the request after %.2fs", 408)
	MsgPolicyEngineRequestInvalid    = ffe("FF21069", "Invalid policy engine request type '%d'")
	MsgPostgresURLMissing            = ffe("FF21070", "URL must be supplied for Postgres persistence")
	MsgPersistenceQueryBuildFailed   = ffe("FF21071", "Persistence query builder failed")
	MsgPersistenceBeginFailed        = ffe("FF21072", "Persistence begin transaction failed")
	MsgPersistenceUpdateFailed       = ffe("FF21073", "Persistence update failed")
	MsgPersistenceCommitFailed       = ffe("FF21074", "Persistence commit failed")
)
