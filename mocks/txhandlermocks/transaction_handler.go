// Code generated by mockery v2.40.1. DO NOT EDIT.

package txhandlermocks

import (
	context "context"

	apitypes "github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"

	ffcapi "github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"

	mock "github.com/stretchr/testify/mock"

	txhandler "github.com/hyperledger/firefly-transaction-manager/pkg/txhandler"
)

// TransactionHandler is an autogenerated mock type for the TransactionHandler type
type TransactionHandler struct {
	mock.Mock
}

// HandleCancelTransaction provides a mock function with given fields: ctx, txID
func (_m *TransactionHandler) HandleCancelTransaction(ctx context.Context, txID string) (*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, txID)

	if len(ret) == 0 {
		panic("no return value specified for HandleCancelTransaction")
	}

	var r0 *apitypes.ManagedTX
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*apitypes.ManagedTX, error)); ok {
		return rf(ctx, txID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleNewContractDeployment provides a mock function with given fields: ctx, txReq
func (_m *TransactionHandler) HandleNewContractDeployment(ctx context.Context, txReq *apitypes.ContractDeployRequest) (*apitypes.ManagedTX, bool, error) {
	ret := _m.Called(ctx, txReq)

	if len(ret) == 0 {
		panic("no return value specified for HandleNewContractDeployment")
	}

	var r0 *apitypes.ManagedTX
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.ContractDeployRequest) (*apitypes.ManagedTX, bool, error)); ok {
		return rf(ctx, txReq)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.ContractDeployRequest) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, txReq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *apitypes.ContractDeployRequest) bool); ok {
		r1 = rf(ctx, txReq)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, *apitypes.ContractDeployRequest) error); ok {
		r2 = rf(ctx, txReq)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HandleNewTransaction provides a mock function with given fields: ctx, txReq
func (_m *TransactionHandler) HandleNewTransaction(ctx context.Context, txReq *apitypes.TransactionRequest) (*apitypes.ManagedTX, bool, error) {
	ret := _m.Called(ctx, txReq)

	if len(ret) == 0 {
		panic("no return value specified for HandleNewTransaction")
	}

	var r0 *apitypes.ManagedTX
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.TransactionRequest) (*apitypes.ManagedTX, bool, error)); ok {
		return rf(ctx, txReq)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.TransactionRequest) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, txReq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *apitypes.TransactionRequest) bool); ok {
		r1 = rf(ctx, txReq)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, *apitypes.TransactionRequest) error); ok {
		r2 = rf(ctx, txReq)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HandleResumeTransaction provides a mock function with given fields: ctx, txID
func (_m *TransactionHandler) HandleResumeTransaction(ctx context.Context, txID string) (*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, txID)

	if len(ret) == 0 {
		panic("no return value specified for HandleResumeTransaction")
	}

	var r0 *apitypes.ManagedTX
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*apitypes.ManagedTX, error)); ok {
		return rf(ctx, txID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleSuspendTransaction provides a mock function with given fields: ctx, txID
func (_m *TransactionHandler) HandleSuspendTransaction(ctx context.Context, txID string) (*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, txID)

	if len(ret) == 0 {
		panic("no return value specified for HandleSuspendTransaction")
	}

	var r0 *apitypes.ManagedTX
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*apitypes.ManagedTX, error)); ok {
		return rf(ctx, txID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleTransactionConfirmations provides a mock function with given fields: ctx, txID, notification
func (_m *TransactionHandler) HandleTransactionConfirmations(ctx context.Context, txID string, notification *apitypes.ConfirmationsNotification) error {
	ret := _m.Called(ctx, txID, notification)

	if len(ret) == 0 {
		panic("no return value specified for HandleTransactionConfirmations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *apitypes.ConfirmationsNotification) error); ok {
		r0 = rf(ctx, txID, notification)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleTransactionReceiptReceived provides a mock function with given fields: ctx, txID, receipt
func (_m *TransactionHandler) HandleTransactionReceiptReceived(ctx context.Context, txID string, receipt *ffcapi.TransactionReceiptResponse) error {
	ret := _m.Called(ctx, txID, receipt)

	if len(ret) == 0 {
		panic("no return value specified for HandleTransactionReceiptReceived")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *ffcapi.TransactionReceiptResponse) error); ok {
		r0 = rf(ctx, txID, receipt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields: ctx, toolkit
func (_m *TransactionHandler) Init(ctx context.Context, toolkit *txhandler.Toolkit) {
	_m.Called(ctx, toolkit)
}

// Start provides a mock function with given fields: ctx
func (_m *TransactionHandler) Start(ctx context.Context) (<-chan struct{}, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 <-chan struct{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (<-chan struct{}, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) <-chan struct{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionHandler creates a new instance of TransactionHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionHandler {
	mock := &TransactionHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
