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

package persistence

import (
	"context"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
)

type EventStreamCheckpoint struct {
	StreamID  *fftypes.UUID                     `json:"streamId"`
	Time      *fftypes.FFTime                   `json:"time"`
	Listeners map[fftypes.UUID]*fftypes.JSONAny `json:"listeners"`
}

type EventStreamPersistence interface {
	StoreCheckpoint(ctx context.Context, checkpoint *EventStreamCheckpoint) error
	ReadCheckpoint(ctx context.Context, streamID *fftypes.UUID) (*EventStreamCheckpoint, error)
	DeleteCheckpoint(ctx context.Context, streamID *fftypes.UUID) error
}
