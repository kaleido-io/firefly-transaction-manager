// Code generated by mockery v1.0.0. DO NOT EDIT.

package eventsmocks

import (
	context "context"

	apitypes "github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Stream is an autogenerated mock type for the Stream type
type Stream struct {
	mock.Mock
}

// AddOrUpdateListener provides a mock function with given fields: ctx, id, updates, reset
func (_m *Stream) AddOrUpdateListener(ctx context.Context, id *fftypes.UUID, updates *apitypes.Listener, reset bool) (*apitypes.Listener, error) {
	ret := _m.Called(ctx, id, updates, reset)

	var r0 *apitypes.Listener
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, *apitypes.Listener, bool) *apitypes.Listener); ok {
		r0 = rf(ctx, id, updates, reset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.Listener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, *apitypes.Listener, bool) error); ok {
		r1 = rf(ctx, id, updates, reset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx
func (_m *Stream) Delete(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveListener provides a mock function with given fields: ctx, id
func (_m *Stream) RemoveListener(ctx context.Context, id *fftypes.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Spec provides a mock function with given fields:
func (_m *Stream) Spec() *apitypes.EventStream {
	ret := _m.Called()

	var r0 *apitypes.EventStream
	if rf, ok := ret.Get(0).(func() *apitypes.EventStream); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.EventStream)
		}
	}

	return r0
}

// Start provides a mock function with given fields: ctx
func (_m *Stream) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Status provides a mock function with given fields:
func (_m *Stream) Status() apitypes.EventStreamStatus {
	ret := _m.Called()

	var r0 apitypes.EventStreamStatus
	if rf, ok := ret.Get(0).(func() apitypes.EventStreamStatus); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(apitypes.EventStreamStatus)
	}

	return r0
}

// Stop provides a mock function with given fields: ctx
func (_m *Stream) Stop(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSpec provides a mock function with given fields: ctx, updates
func (_m *Stream) UpdateSpec(ctx context.Context, updates *apitypes.EventStream) error {
	ret := _m.Called(ctx, updates)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.EventStream) error); ok {
		r0 = rf(ctx, updates)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
