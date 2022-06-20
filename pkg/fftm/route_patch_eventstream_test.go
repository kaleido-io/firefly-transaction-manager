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

package fftm

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/stretchr/testify/assert"
)

func TestPatchEventStream(t *testing.T) {

	url, m, done := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {})
	defer done()

	err := m.Start()
	assert.NoError(t, err)

	// Create stream
	var es apitypes.EventStream
	res, err := resty.New().R().
		SetBody(&apitypes.EventStream{
			Name: strPtr("my event stream"),
		}).
		SetResult(&es).
		Post(url + "/eventstreams")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode())

	// Then update it
	res, err = resty.New().R().
		SetBody(&apitypes.EventStream{
			Name: strPtr("my renamed event stream"),
		}).
		SetResult(&es).
		Patch(url + "/eventstreams/" + es.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode())
	assert.NotNil(t, es.ID)
	assert.NotNil(t, es.Created)
	assert.NotEqual(t, es.Created, es.Updated)
	assert.Equal(t, "my renamed event stream", *es.Name)

}