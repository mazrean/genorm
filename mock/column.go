package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	genorm "github.com/mazrean/genorm"
)

// MockTypedTableColumn is a mock of Column interface.
type MockTypedTableColumn[T genorm.Table, S genorm.ExprType] struct {
	ctrl     *gomock.Controller
	recorder *MockTypedTableColumnMockRecorder[T, S]
}

// MockTypedTableColumnMockRecorder is the mock recorder for MockTypedTableColumn.
type MockTypedTableColumnMockRecorder[T genorm.Table, S genorm.ExprType] struct {
	mock *MockTypedTableColumn[T, S]
}

// NewMockTypedTableColumn creates a new mock instance.
func NewMockTypedTableColumn[T genorm.Table, S genorm.ExprType](ctrl *gomock.Controller) *MockTypedTableColumn[T, S] {
	mock := &MockTypedTableColumn[T, S]{ctrl: ctrl}
	mock.recorder = &MockTypedTableColumnMockRecorder[T, S]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTypedTableColumn[T, S]) EXPECT() *MockTypedTableColumnMockRecorder[T, S] {
	return m.recorder
}

// ColumnName mocks base method.
func (m *MockTypedTableColumn[_, _]) ColumnName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ColumnName")
	ret0, _ := ret[0].(string)
	return ret0
}

// ColumnName indicates an expected call of ColumnName.
func (mr *MockTypedTableColumnMockRecorder[T, S]) ColumnName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ColumnName", reflect.TypeOf((*MockTypedTableColumn[T, S])(nil).ColumnName))
}

// Expr mocks base method.
func (m *MockTypedTableColumn[_, _]) Expr() (string, []genorm.ExprType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expr")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]genorm.ExprType)
	ret2, _ := ret[2].([]error)
	return ret0, ret1, ret2
}

// Expr indicates an expected call of Expr.
func (mr *MockTypedTableColumnMockRecorder[T, S]) Expr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expr", reflect.TypeOf((*MockTypedTableColumn[T, S])(nil).Expr))
}

// TableExpr mocks base method.
func (m *MockTypedTableColumn[T, _]) TableExpr(t T) (string, []genorm.ExprType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableExpr", t)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]genorm.ExprType)
	ret2, _ := ret[2].([]error)
	return ret0, ret1, ret2
}

// TableExpr indicates an expected call of Expr.
func (mr *MockTypedTableColumnMockRecorder[T, S]) TableExpr(t T) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableExpr", reflect.TypeOf((*MockTypedTableColumn[T, S])(nil).Expr), t)
}

// TypedExpr mocks base method.
func (m *MockTypedTableColumn[_, S]) TypedExpr(s S) (string, []genorm.ExprType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TypedExpr", s)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]genorm.ExprType)
	ret2, _ := ret[2].([]error)
	return ret0, ret1, ret2
}

// Expr indicates an expected call of Expr.
func (mr *MockTypedTableColumnMockRecorder[T, S]) TypedExpr(s S) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TypedExpr", reflect.TypeOf((*MockTypedTableColumn[T, S])(nil).Expr), s)
}

// SQLColumnName mocks base method.
func (m *MockTypedTableColumn[_, _]) SQLColumnName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SQLColumnName")
	ret0, _ := ret[0].(string)
	return ret0
}

// SQLColumnName indicates an expected call of SQLColumnName.
func (mr *MockTypedTableColumnMockRecorder[T, S]) SQLColumnName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SQLColumnName", reflect.TypeOf((*MockTypedTableColumn[T, S])(nil).SQLColumnName))
}

// TableName mocks base method.
func (m *MockTypedTableColumn[_, _]) TableName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableName")
	ret0, _ := ret[0].(string)
	return ret0
}

// TableName indicates an expected call of TableName.
func (mr *MockTypedTableColumnMockRecorder[T, S]) TableName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableName", reflect.TypeOf((*MockTypedTableColumn[T, S])(nil).TableName))
}
