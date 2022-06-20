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
	"fmt"
	"net/http"
	"testing"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-transaction-manager/mocks/ffcapimocks"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRestoreStreamsAndListenersOK(t *testing.T) {

	_, m, done := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {})
	defer done()

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything, mock.Anything).Return("EventSig(uint256)", *fftypes.JSONAnyPtr(`{}`), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerAddResponse{}, ffcapi.ErrorReason(""), nil)
	mfc.On("EventListenerRemove", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerRemoveResponse{}, ffcapi.ErrorReason(""), nil).Maybe()

	falsy := false

	es1 := &apitypes.EventStream{ID: apitypes.UUIDVersion1(), Name: strPtr("stream1"), Suspended: &falsy}
	err := m.persistence.WriteStream(m.ctx, es1)
	assert.NoError(t, err)

	e1l1 := &apitypes.Listener{ID: apitypes.UUIDVersion1(), Name: "listener1", StreamID: es1.ID}
	err = m.persistence.WriteListener(m.ctx, e1l1)
	assert.NoError(t, err)

	e1l2 := &apitypes.Listener{ID: apitypes.UUIDVersion1(), Name: "listener2", StreamID: es1.ID}
	err = m.persistence.WriteListener(m.ctx, e1l2)
	assert.NoError(t, err)

	es2 := &apitypes.EventStream{ID: apitypes.UUIDVersion1(), Name: strPtr("stream2"), Suspended: &falsy}
	err = m.persistence.WriteStream(m.ctx, es2)
	assert.NoError(t, err)

	e2l1 := &apitypes.Listener{ID: apitypes.UUIDVersion1(), Name: "listener3", StreamID: es2.ID}
	err = m.persistence.WriteListener(m.ctx, e2l1)
	assert.NoError(t, err)

	err = m.Start()
	assert.NoError(t, err)

	assert.Equal(t, es1.ID, m.streamsByName["stream1"])
	assert.Equal(t, es2.ID, m.streamsByName["stream2"])

	mfc.AssertExpectations(t)

}

func TestRestoreStreamsReadFailed(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	mp.On("ListStreams", m.ctx, (*fftypes.UUID)(nil), startupPaginationLimit).Return(nil, fmt.Errorf("pop"))

	err := m.restoreStreams()
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestRestoreListenersReadFailed(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	mp.On("ListListeners", m.ctx, (*fftypes.UUID)(nil), startupPaginationLimit).Return(nil, fmt.Errorf("pop"))

	err := m.restoreListeners()
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestRestoreStreamsValidateFail(t *testing.T) {

	_, m, done := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {})
	defer done()

	falsy := false
	es1 := &apitypes.EventStream{ID: apitypes.UUIDVersion1(), Name: strPtr(""), Suspended: &falsy}
	err := m.persistence.WriteStream(m.ctx, es1)
	assert.NoError(t, err)

	err = m.restoreStreams()
	assert.Regexp(t, "FF21028", err)

}

func TestRestoreListenersStartFail(t *testing.T) {

	_, m, done := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {})
	defer done()

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything, mock.Anything).Return("EventSig(uint256)", *fftypes.JSONAnyPtr(`{}`), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), fmt.Errorf("pop"))

	falsy := false
	es1 := &apitypes.EventStream{ID: apitypes.UUIDVersion1(), Name: strPtr("stream1"), Suspended: &falsy}
	err := m.persistence.WriteStream(m.ctx, es1)
	assert.NoError(t, err)

	e1l1 := &apitypes.Listener{ID: apitypes.UUIDVersion1(), Name: "listener1", StreamID: es1.ID}
	err = m.persistence.WriteListener(m.ctx, e1l1)
	assert.NoError(t, err)

	err = m.restoreStreams()
	assert.Regexp(t, "pop", err)

	mfc.AssertExpectations(t)

}

func TestDeleteStartedListener(t *testing.T) {

	_, m, done := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {})
	defer done()

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything, mock.Anything).Return("EventSig(uint256)", *fftypes.JSONAnyPtr(`{}`), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), nil)
	mfc.On("EventListenerRemove", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerRemoveResponse{}, ffcapi.ErrorReason(""), nil)

	falsy := false
	es1 := &apitypes.EventStream{ID: apitypes.UUIDVersion1(), Name: strPtr("stream1"), Suspended: &falsy}
	err := m.persistence.WriteStream(m.ctx, es1)
	assert.NoError(t, err)

	e1l1 := &apitypes.Listener{ID: apitypes.UUIDVersion1(), Name: "listener1", StreamID: es1.ID}
	err = m.persistence.WriteListener(m.ctx, e1l1)
	assert.NoError(t, err)

	err = m.Start()
	assert.NoError(t, err)

	err = m.deleteStream(m.ctx, es1.ID.String())
	assert.NoError(t, err)

	mfc.AssertExpectations(t)

}

func TestDeleteStartedListenerFail(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	esID := apitypes.UUIDVersion1()
	lID := apitypes.UUIDVersion1()
	mp.On("ListStreamListeners", m.ctx, (*fftypes.UUID)(nil), startupPaginationLimit, esID).Return([]*apitypes.Listener{
		{ID: lID, StreamID: esID},
	}, nil)
	mp.On("DeleteListener", m.ctx, lID).Return(fmt.Errorf("pop"))

	err := m.deleteAllStreamListeners(m.ctx, esID)
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestDeleteStreamBadID(t *testing.T) {

	_, _, m := newMockPersistenceManager(t)

	err := m.deleteStream(m.ctx, "Bad ID")
	assert.Regexp(t, "FF00138", err)

}

func TestDeleteStreamListenerPersistenceFail(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	esID := apitypes.UUIDVersion1()
	mp.On("ListStreamListeners", m.ctx, (*fftypes.UUID)(nil), startupPaginationLimit, esID).Return(nil, fmt.Errorf("pop"))

	err := m.deleteStream(m.ctx, esID.String())
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestDeleteStreamPersistenceFail(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	esID := apitypes.UUIDVersion1()
	mp.On("ListStreamListeners", m.ctx, (*fftypes.UUID)(nil), startupPaginationLimit, esID).Return([]*apitypes.Listener{}, nil)
	mp.On("DeleteStream", m.ctx, esID).Return(fmt.Errorf("pop"))

	err := m.deleteStream(m.ctx, esID.String())
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestDeleteStreamNotInitialized(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	esID := apitypes.UUIDVersion1()
	mp.On("ListStreamListeners", m.ctx, (*fftypes.UUID)(nil), startupPaginationLimit, esID).Return([]*apitypes.Listener{}, nil)
	mp.On("DeleteStream", m.ctx, esID).Return(nil)

	err := m.deleteStream(m.ctx, esID.String())
	assert.NoError(t, err)

	mp.AssertExpectations(t)
}

func TestCreateRenameStreamNameReservation(t *testing.T) {

	mp, _, m := newMockPersistenceManager(t)

	mp.On("WriteStream", m.ctx, mock.Anything).Return(fmt.Errorf("temporary")).Once()
	mp.On("DeleteCheckpoint", m.ctx, mock.Anything).Return(fmt.Errorf("temporary")).Once()
	mp.On("WriteStream", m.ctx, mock.Anything).Return(nil)

	// Reject missing name
	_, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{})
	assert.Regexp(t, "FF21028", err)

	// Attempt to start and encounter a temporary error
	_, err = m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("Name1")})
	assert.Regexp(t, "temporary", err)

	// Ensure we still allow use of the name after the glitch is fixed
	es1, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("Name1")})
	assert.NoError(t, err)

	// Ensure we can't create another stream of same name
	_, err = m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("Name1")})
	assert.Regexp(t, "FF21047", err)

	// Create a second stream to test clash on rename
	es2, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("Name2")})
	assert.NoError(t, err)

	// Check for clash
	_, err = m.updateStream(m.ctx, es1.ID.String(), &apitypes.EventStream{Name: strPtr("Name2")})
	assert.Regexp(t, "FF21047", err)

	// Check for no-op rename to self
	_, err = m.updateStream(m.ctx, es2.ID.String(), &apitypes.EventStream{Name: strPtr("Name2")})
	assert.NoError(t, err)

	mp.AssertExpectations(t)
}

func TestCreateStreamValidateFail(t *testing.T) {

	_, _, m := newMockPersistenceManager(t)

	wrongType := apitypes.DistributionMode("wrong")
	_, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("stream1"), Type: &wrongType})
	assert.Regexp(t, "FF21029", err)

}

func TestCreateAndStoreNewStreamListenerBadID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.createAndStoreNewStreamListener(m.ctx, "bad", nil)
	assert.Regexp(t, "FF00138", err)
}

func TestUpdateExistingListenerNotFound(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("GetListener", m.ctx, mock.Anything).Return(nil, nil)

	_, err := m.updateExistingListener(m.ctx, apitypes.UUIDVersion1().String(), apitypes.UUIDVersion1().String(), &apitypes.Listener{})
	assert.Regexp(t, "FF21046", err)

	mp.AssertExpectations(t)
}

func TestCreateOrUpdateListenerNotFound(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.createOrUpdateListener(m.ctx, &apitypes.Listener{ID: apitypes.UUIDVersion1(), StreamID: apitypes.UUIDVersion1()})
	assert.Regexp(t, "FF21045", err)

}

func TestCreateOrUpdateListenerFail(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("WriteStream", m.ctx, mock.Anything).Return(nil)

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything, mock.Anything).Return("EventSig(uint256)", *fftypes.JSONAnyPtr(`{}`), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), fmt.Errorf("pop"))

	es, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("stream1")})

	_, err = m.createOrUpdateListener(m.ctx, &apitypes.Listener{ID: apitypes.UUIDVersion1(), StreamID: es.ID})
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestCreateOrUpdateListenerWriteFail(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("WriteStream", m.ctx, mock.Anything).Return(nil)
	mp.On("WriteListener", m.ctx, mock.Anything).Return(fmt.Errorf("pop"))

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything, mock.Anything).Return("EventSig(uint256)", *fftypes.JSONAnyPtr(`{}`), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), nil)
	mfc.On("EventListenerRemove", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerRemoveResponse{}, ffcapi.ErrorReason(""), nil)

	es, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("stream1")})

	_, err = m.createOrUpdateListener(m.ctx, &apitypes.Listener{ID: apitypes.UUIDVersion1(), StreamID: es.ID})
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)
}

func TestDeleteListenerBadID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	err := m.deleteListener(m.ctx, "bad ID", "bad ID")
	assert.Regexp(t, "FF00138", err)

}

func TestDeleteListenerStreamNotFound(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	l1 := &apitypes.Listener{ID: apitypes.UUIDVersion1(), StreamID: apitypes.UUIDVersion1()}
	mp.On("GetListener", m.ctx, mock.Anything).Return(l1, nil)

	err := m.deleteListener(m.ctx, l1.StreamID.String(), l1.ID.String())
	assert.Regexp(t, "FF21045", err)

	mp.AssertExpectations(t)

}

func TestDeleteListenerFail(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("WriteStream", m.ctx, mock.Anything).Return(nil)
	mp.On("WriteListener", m.ctx, mock.Anything).Return(nil)

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything, mock.Anything).Return("EventSig(uint256)", *fftypes.JSONAnyPtr(`{}`), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), nil)
	mfc.On("EventListenerRemove", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), fmt.Errorf("pop"))

	es, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("stream1")})

	l1, err := m.createOrUpdateListener(m.ctx, &apitypes.Listener{ID: apitypes.UUIDVersion1(), StreamID: es.ID})
	assert.NoError(t, err)

	mp.On("GetListener", m.ctx, mock.Anything).Return(l1, nil)

	err = m.deleteListener(m.ctx, l1.StreamID.String(), l1.ID.String())
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)

}

func TestUpdateStreamBadID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.updateStream(m.ctx, "bad ID", &apitypes.EventStream{})
	assert.Regexp(t, "FF00138", err)

}

func TestUpdateStreamNotFound(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.updateStream(m.ctx, apitypes.UUIDVersion1().String(), &apitypes.EventStream{})
	assert.Regexp(t, "FF21045", err)

}

func TestUpdateStreamBadChanges(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("WriteStream", m.ctx, mock.Anything).Return(nil)

	es, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("stream1")})

	wrongType := apitypes.DistributionMode("wrong")
	_, err = m.updateStream(m.ctx, es.ID.String(), &apitypes.EventStream{Type: &wrongType})
	assert.Regexp(t, "FF21029", err)

}

func TestUpdateStreamWriteFail(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("WriteStream", m.ctx, mock.Anything).Return(nil).Once()
	mp.On("WriteStream", m.ctx, mock.Anything).Return(fmt.Errorf("pop"))

	es, err := m.createAndStoreNewStream(m.ctx, &apitypes.EventStream{Name: strPtr("stream1")})

	_, err = m.updateStream(m.ctx, es.ID.String(), &apitypes.EventStream{})
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)

}

func TestGetStreamBadID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getStream(m.ctx, "bad ID")
	assert.Regexp(t, "FF00138", err)

}

func TestGetStreamNotFound(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getStream(m.ctx, apitypes.UUIDVersion1().String())
	assert.Regexp(t, "FF21045", err)

}

func TestGetStreamsBadLimit(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getStreams(m.ctx, "", "wrong")
	assert.Regexp(t, "FF21044", err)

}

func TestGetListenerBadAfter(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getListeners(m.ctx, "!bad UUID", "")
	assert.Regexp(t, "FF00138", err)

}

func TestGetListenerBadStreamID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getListener(m.ctx, "bad ID", apitypes.UUIDVersion1().String())
	assert.Regexp(t, "FF00138", err)

}

func TestGetListenerBadListenerID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getListener(m.ctx, apitypes.UUIDVersion1().String(), "bad ID")
	assert.Regexp(t, "FF00138", err)

}

func TestGetListenerLookupErr(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("GetListener", m.ctx, mock.Anything).Return(nil, fmt.Errorf("pop"))

	_, err := m.getListener(m.ctx, apitypes.UUIDVersion1().String(), apitypes.UUIDVersion1().String())
	assert.Regexp(t, "pop", err)

	mp.AssertExpectations(t)

}

func TestGetListenerNotFound(t *testing.T) {
	mp, _, m := newMockPersistenceManager(t)

	mp.On("GetListener", m.ctx, mock.Anything).Return(nil, nil)

	_, err := m.getListener(m.ctx, apitypes.UUIDVersion1().String(), apitypes.UUIDVersion1().String())
	assert.Regexp(t, "FF21046", err)

	mp.AssertExpectations(t)

}

func TestGetStreamListenersBadLimit(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getStreamListeners(m.ctx, "", "!bad limit", apitypes.UUIDVersion1().String())
	assert.Regexp(t, "FF21044", err)

}

func TestGetStreamListenersBadStreamID(t *testing.T) {
	_, _, m := newMockPersistenceManager(t)

	_, err := m.getStreamListeners(m.ctx, "", "", "bad ID")
	assert.Regexp(t, "FF00138", err)

}