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

package apitypes

import (
	"encoding/json"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
)

var (
	EventStreamTypeWebhook   = fftypes.FFEnumValue("estype", "webhook")
	EventStreamTypeWebSocket = fftypes.FFEnumValue("estype", "websocket")
)

type EventStream struct {
	ID        *fftypes.UUID    `ffstruct:"eventstream" json:"id"`
	Created   *fftypes.FFTime  `ffstruct:"eventstream" json:"created"`
	Updated   *fftypes.FFTime  `ffstruct:"eventstream" json:"updated"`
	Name      *string          `ffstruct:"eventstream" json:"name,omitempty"`
	Suspended *bool            `ffstruct:"eventstream" json:"suspended,omitempty"`
	Type      *EventStreamType `ffstruct:"eventstream" json:"type,omitempty" ffenum:"estype"`

	ErrorHandling     *ErrorHandlingType  `ffstruct:"eventstream" json:"errorHandling"`
	BatchSize         *uint64             `ffstruct:"eventstream" json:"batchSize"`
	BatchTimeout      *fftypes.FFDuration `ffstruct:"eventstream" json:"batchTimeout"`
	RetryTimeout      *fftypes.FFDuration `ffstruct:"eventstream" json:"retryTimeout"`
	BlockedRetryDelay *fftypes.FFDuration `ffstruct:"eventstream" json:"blockedRetryDelay"`

	EthCompatBatchTimeoutMS       *uint64 `ffstruct:"eventstream" json:"batchTimeoutMS,omitempty"`       // input only, for backwards compatibility
	EthCompatRetryTimeoutSec      *uint64 `ffstruct:"eventstream" json:"retryTimeoutSec,omitempty"`      // input only, for backwards compatibility
	EthCompatBlockedRetryDelaySec *uint64 `ffstruct:"eventstream" json:"blockedRetryDelaySec,omitempty"` // input only, for backwards compatibility

	Webhook   *WebhookConfig   `ffstruct:"eventstream" json:"webhook,omitempty"`
	WebSocket *WebSocketConfig `ffstruct:"eventstream" json:"websocket,omitempty"`
}

type EventStreamStatus string

const (
	EventStreamStatusStarted  EventStreamStatus = "started"
	EventStreamStatusStopping EventStreamStatus = "stopping"
	EventStreamStatusStopped  EventStreamStatus = "stopped"
	EventStreamStatusDeleted  EventStreamStatus = "deleted"
)

type EventStreamWithStatus struct {
	EventStream
	Status EventStreamStatus `ffstruct:"eventstream" json:"status"`
}

type EventStreamCheckpoint struct {
	StreamID  *fftypes.UUID                    `json:"streamId"`
	Time      *fftypes.FFTime                  `json:"time"`
	Listeners map[fftypes.UUID]json.RawMessage `json:"listeners"`
}
