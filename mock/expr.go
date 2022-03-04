package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	genorm "github.com/mazrean/genorm"
)

// MockTypedTableExpr is a mock of Expr interface.
type MockTypedTableExpr[T genorm.Table, S genorm.ExprType] struct {
	ctrl     *gomock.Controller
	recorder *MockTypedTableExprMockRecorder[T, S]
}

// MockTypedTableExprMockRecorder is the mock recorder for MockTypedTableExpr.
type MockTypedTableExprMockRecorder[T genorm.Table, S genorm.ExprType] struct {
	mock *MockTypedTableExpr[T, S]
}

// NewMockTypedTableExpr creates a new mock instance.
func NewMockTypedTableExpr[T genorm.Table, S genorm.ExprType](ctrl *gomock.Controller) *MockTypedTableExpr[T, S] {
	mock := &MockTypedTableExpr[T, S]{ctrl: ctrl}
	mock.recorder = &MockTypedTableExprMockRecorder[T, S]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTypedTableExpr[T, S]) EXPECT() *MockTypedTableExprMockRecorder[T, S] {
	return m.recorder
}

// Expr mocks base method.
func (m *MockTypedTableExpr[_, _]) Expr() (string, []genorm.ExprType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expr")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]genorm.ExprType)
	ret2, _ := ret[2].([]error)
	return ret0, ret1, ret2
}

// Expr indicates an expected call of Expr.
func (mr *MockTypedTableExprMockRecorder[T, S]) Expr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expr", reflect.TypeOf((*MockTypedTableExpr[T, S])(nil).Expr))
}

// TableExpr mocks base method.
func (m *MockTypedTableExpr[T, _]) TableExpr(t T) (string, []genorm.ExprType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableExpr", t)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]genorm.ExprType)
	ret2, _ := ret[2].([]error)
	return ret0, ret1, ret2
}

// TableExpr indicates an expected call of Expr.
func (mr *MockTypedTableExprMockRecorder[T, S]) TableExpr(t T) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableExpr", reflect.TypeOf((*MockTypedTableExpr[T, S])(nil).Expr), t)
}

// TypedExpr mocks base method.
func (m *MockTypedTableExpr[_, S]) TypedExpr(s S) (string, []genorm.ExprType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TypedExpr", s)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]genorm.ExprType)
	ret2, _ := ret[2].([]error)
	return ret0, ret1, ret2
}

// TableExpr indicates an expected call of Expr.
func (mr *MockTypedTableExprMockRecorder[T, S]) TypedExpr(s S) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TypedExpr", reflect.TypeOf((*MockTypedTableExpr[T, S])(nil).Expr), s)
}
