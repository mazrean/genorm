package genorm_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mazrean/genorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserTable is a test table struct for Select binding tests
type TestUserTable struct {
	ID        genorm.WrappedPrimitive[int64]
	Name      genorm.WrappedPrimitive[string]
	Age       genorm.WrappedPrimitive[int32]
	Score     genorm.WrappedPrimitive[float64]
	Active    genorm.WrappedPrimitive[bool]
	CreatedAt genorm.WrappedPrimitive[time.Time]
}

func (t *TestUserTable) TableName() string {
	return "users"
}

func (t *TestUserTable) Columns() []genorm.Column {
	return []genorm.Column{
		&testColumn{table: "users", column: "id", sqlColumn: "users.id"},
		&testColumn{table: "users", column: "name", sqlColumn: "users.name"},
		&testColumn{table: "users", column: "age", sqlColumn: "users.age"},
		&testColumn{table: "users", column: "score", sqlColumn: "users.score"},
		&testColumn{table: "users", column: "active", sqlColumn: "users.active"},
		&testColumn{table: "users", column: "created_at", sqlColumn: "users.created_at"},
	}
}

func (t *TestUserTable) ColumnMap() map[string]genorm.ColumnFieldExprType {
	return map[string]genorm.ColumnFieldExprType{
		"users.id":         &t.ID,
		"users.name":       &t.Name,
		"users.age":        &t.Age,
		"users.score":      &t.Score,
		"users.active":     &t.Active,
		"users.created_at": &t.CreatedAt,
	}
}

func (t *TestUserTable) Expr() (string, []genorm.ExprType, []error) {
	return "users", nil, nil
}

func (t *TestUserTable) GetErrors() []error {
	return nil
}

// testColumn is a simple column implementation for testing
type testColumn struct {
	table     string
	column    string
	sqlColumn string
}

func (c *testColumn) TableName() string {
	return c.table
}

func (c *testColumn) ColumnName() string {
	return c.column
}

func (c *testColumn) SQLColumnName() string {
	return c.sqlColumn
}

func (c *testColumn) Expr() (string, []genorm.ExprType, []error) {
	return c.sqlColumn, nil, nil
}

func TestSelectGetAll_SingleRow(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	}).AddRow(
		int64(1), "Alice", int32(30), float64(95.5), true, createdAt,
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestUserTable{}).GetAllCtx(ctx, db)

	require.NoError(t, err)
	require.Len(t, results, 1)

	// Verify the values are correctly bound
	id, valid := results[0].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(1), id)

	name, valid := results[0].Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Alice", name)

	age, valid := results[0].Age.Val()
	assert.True(t, valid)
	assert.Equal(t, int32(30), age)

	score, valid := results[0].Score.Val()
	assert.True(t, valid)
	assert.Equal(t, 95.5, score)

	active, valid := results[0].Active.Val()
	assert.True(t, valid)
	assert.True(t, active)

	createdAtVal, valid := results[0].CreatedAt.Val()
	assert.True(t, valid)
	assert.Equal(t, createdAt, createdAtVal)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectGetAll_MultipleRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	createdAt1 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	createdAt2 := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
	createdAt3 := time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	}).AddRow(
		int64(1), "Alice", int32(30), float64(95.5), true, createdAt1,
	).AddRow(
		int64(2), "Bob", int32(25), float64(88.0), false, createdAt2,
	).AddRow(
		int64(3), "Charlie", int32(35), float64(92.3), true, createdAt3,
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestUserTable{}).GetAllCtx(ctx, db)

	require.NoError(t, err)
	require.Len(t, results, 3)

	// Verify first row
	id, valid := results[0].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(1), id)
	name, valid := results[0].Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Alice", name)

	// Verify second row
	id, valid = results[1].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(2), id)
	name, valid = results[1].Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Bob", name)

	// Verify third row
	id, valid = results[2].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(3), id)
	name, valid = results[2].Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Charlie", name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectGetAll_EmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	})

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestUserTable{}).GetAllCtx(ctx, db)

	require.NoError(t, err)
	assert.Len(t, results, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectGetAll_NullValues(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	}).AddRow(
		int64(1), sql.NullString{Valid: false}, sql.NullInt32{Valid: false}, sql.NullFloat64{Valid: false}, sql.NullBool{Valid: false}, sql.NullTime{Valid: false},
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestUserTable{}).GetAllCtx(ctx, db)

	require.NoError(t, err)
	require.Len(t, results, 1)

	// ID should be valid
	id, valid := results[0].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(1), id)

	// Other fields should be invalid (null)
	_, valid = results[0].Name.Val()
	assert.False(t, valid)

	_, valid = results[0].Age.Val()
	assert.False(t, valid)

	_, valid = results[0].Score.Val()
	assert.False(t, valid)

	_, valid = results[0].Active.Val()
	assert.False(t, valid)

	_, valid = results[0].CreatedAt.Val()
	assert.False(t, valid)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectGet_SingleRow(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	}).AddRow(
		int64(1), "Alice", int32(30), float64(95.5), true, createdAt,
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users LIMIT 1`).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := genorm.Select(&TestUserTable{}).GetCtx(ctx, db)

	require.NoError(t, err)

	// Verify the values are correctly bound
	id, valid := result.ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(1), id)

	name, valid := result.Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Alice", name)

	age, valid := result.Age.Val()
	assert.True(t, valid)
	assert.Equal(t, int32(30), age)

	score, valid := result.Score.Val()
	assert.True(t, valid)
	assert.Equal(t, 95.5, score)

	active, valid := result.Active.Val()
	assert.True(t, valid)
	assert.True(t, active)

	createdAtVal, valid := result.CreatedAt.Val()
	assert.True(t, valid)
	assert.Equal(t, createdAt, createdAtVal)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectGet_NoRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	})

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users LIMIT 1`).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := genorm.Select(&TestUserTable{}).GetCtx(ctx, db)

	assert.Error(t, err)
	assert.ErrorIs(t, err, genorm.ErrRecordNotFound)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectGetAll_ScanError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Return wrong type for age field (string instead of int)
	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	}).AddRow(
		int64(1), "Alice", "invalid_age", float64(95.5), true, time.Now(),
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestUserTable{}).GetAllCtx(ctx, db)

	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestSelectGet_ScanError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Return wrong type for age field (string instead of int)
	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0", "users_age_0", "users_score_0", "users_active_0", "users_created_at_0",
	}).AddRow(
		int64(1), "Alice", "invalid_age", float64(95.5), true, time.Now(),
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0, users\.age AS users_age_0, users\.score AS users_score_0, users\.active AS users_active_0, users\.created_at AS users_created_at_0 FROM users LIMIT 1`).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := genorm.Select(&TestUserTable{}).GetCtx(ctx, db)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestPartialUserTable is a test table struct for partial column selection tests
type TestPartialUserTable struct {
	ID   genorm.WrappedPrimitive[int64]
	Name genorm.WrappedPrimitive[string]
}

func (t *TestPartialUserTable) TableName() string {
	return "users"
}

func (t *TestPartialUserTable) Columns() []genorm.Column {
	return []genorm.Column{
		&testColumn{table: "users", column: "id", sqlColumn: "users.id"},
		&testColumn{table: "users", column: "name", sqlColumn: "users.name"},
	}
}

func (t *TestPartialUserTable) ColumnMap() map[string]genorm.ColumnFieldExprType {
	return map[string]genorm.ColumnFieldExprType{
		"users.id":   &t.ID,
		"users.name": &t.Name,
	}
}

func (t *TestPartialUserTable) Expr() (string, []genorm.ExprType, []error) {
	return "users", nil, nil
}

func (t *TestPartialUserTable) GetErrors() []error {
	return nil
}

func TestSelectGetAll_PartialColumns(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"users_id_0", "users_name_0",
	}).AddRow(
		int64(1), "Alice",
	).AddRow(
		int64(2), "Bob",
	)

	mock.ExpectQuery(`SELECT users\.id AS users_id_0, users\.name AS users_name_0 FROM users`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestPartialUserTable{}).GetAllCtx(ctx, db)

	require.NoError(t, err)
	require.Len(t, results, 2)

	// Verify first row
	id, valid := results[0].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(1), id)
	name, valid := results[0].Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Alice", name)

	// Verify second row
	id, valid = results[1].ID.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(2), id)
	name, valid = results[1].Name.Val()
	assert.True(t, valid)
	assert.Equal(t, "Bob", name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestNumericTypes tests various numeric type bindings
type TestNumericTable struct {
	Int8Val   genorm.WrappedPrimitive[int8]
	Int16Val  genorm.WrappedPrimitive[int16]
	Int32Val  genorm.WrappedPrimitive[int32]
	Int64Val  genorm.WrappedPrimitive[int64]
	Uint8Val  genorm.WrappedPrimitive[uint8]
	Uint16Val genorm.WrappedPrimitive[uint16]
	Uint32Val genorm.WrappedPrimitive[uint32]
	Uint64Val genorm.WrappedPrimitive[uint64]
	Float32   genorm.WrappedPrimitive[float32]
	Float64   genorm.WrappedPrimitive[float64]
}

func (t *TestNumericTable) TableName() string {
	return "numeric_test"
}

func (t *TestNumericTable) Columns() []genorm.Column {
	return []genorm.Column{
		&testColumn{table: "numeric_test", column: "int8_val", sqlColumn: "numeric_test.int8_val"},
		&testColumn{table: "numeric_test", column: "int16_val", sqlColumn: "numeric_test.int16_val"},
		&testColumn{table: "numeric_test", column: "int32_val", sqlColumn: "numeric_test.int32_val"},
		&testColumn{table: "numeric_test", column: "int64_val", sqlColumn: "numeric_test.int64_val"},
		&testColumn{table: "numeric_test", column: "uint8_val", sqlColumn: "numeric_test.uint8_val"},
		&testColumn{table: "numeric_test", column: "uint16_val", sqlColumn: "numeric_test.uint16_val"},
		&testColumn{table: "numeric_test", column: "uint32_val", sqlColumn: "numeric_test.uint32_val"},
		&testColumn{table: "numeric_test", column: "uint64_val", sqlColumn: "numeric_test.uint64_val"},
		&testColumn{table: "numeric_test", column: "float32", sqlColumn: "numeric_test.float32"},
		&testColumn{table: "numeric_test", column: "float64", sqlColumn: "numeric_test.float64"},
	}
}

func (t *TestNumericTable) ColumnMap() map[string]genorm.ColumnFieldExprType {
	return map[string]genorm.ColumnFieldExprType{
		"numeric_test.int8_val":   &t.Int8Val,
		"numeric_test.int16_val":  &t.Int16Val,
		"numeric_test.int32_val":  &t.Int32Val,
		"numeric_test.int64_val":  &t.Int64Val,
		"numeric_test.uint8_val":  &t.Uint8Val,
		"numeric_test.uint16_val": &t.Uint16Val,
		"numeric_test.uint32_val": &t.Uint32Val,
		"numeric_test.uint64_val": &t.Uint64Val,
		"numeric_test.float32":    &t.Float32,
		"numeric_test.float64":    &t.Float64,
	}
}

func (t *TestNumericTable) Expr() (string, []genorm.ExprType, []error) {
	return "numeric_test", nil, nil
}

func (t *TestNumericTable) GetErrors() []error {
	return nil
}

func TestSelectGetAll_NumericTypes(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"numeric_test_int8_val_0", "numeric_test_int16_val_0", "numeric_test_int32_val_0", "numeric_test_int64_val_0",
		"numeric_test_uint8_val_0", "numeric_test_uint16_val_0", "numeric_test_uint32_val_0", "numeric_test_uint64_val_0",
		"numeric_test_float32_0", "numeric_test_float64_0",
	}).AddRow(
		int8(127), int16(32767), int32(2147483647), int64(9223372036854775807),
		uint8(255), uint16(65535), uint32(4294967295), uint64(1234567890),
		float32(3.14), float64(2.718281828),
	)

	mock.ExpectQuery(`SELECT numeric_test\.int8_val AS numeric_test_int8_val_0, numeric_test\.int16_val AS numeric_test_int16_val_0, numeric_test\.int32_val AS numeric_test_int32_val_0, numeric_test\.int64_val AS numeric_test_int64_val_0, numeric_test\.uint8_val AS numeric_test_uint8_val_0, numeric_test\.uint16_val AS numeric_test_uint16_val_0, numeric_test\.uint32_val AS numeric_test_uint32_val_0, numeric_test\.uint64_val AS numeric_test_uint64_val_0, numeric_test\.float32 AS numeric_test_float32_0, numeric_test\.float64 AS numeric_test_float64_0 FROM numeric_test`).
		WillReturnRows(rows)

	ctx := context.Background()
	results, err := genorm.Select(&TestNumericTable{}).GetAllCtx(ctx, db)

	require.NoError(t, err)
	require.Len(t, results, 1)

	// Verify all numeric types are correctly bound
	int8Val, valid := results[0].Int8Val.Val()
	assert.True(t, valid)
	assert.Equal(t, int8(127), int8Val)

	int16Val, valid := results[0].Int16Val.Val()
	assert.True(t, valid)
	assert.Equal(t, int16(32767), int16Val)

	int32Val, valid := results[0].Int32Val.Val()
	assert.True(t, valid)
	assert.Equal(t, int32(2147483647), int32Val)

	int64Val, valid := results[0].Int64Val.Val()
	assert.True(t, valid)
	assert.Equal(t, int64(9223372036854775807), int64Val)

	uint8Val, valid := results[0].Uint8Val.Val()
	assert.True(t, valid)
	assert.Equal(t, uint8(255), uint8Val)

	uint16Val, valid := results[0].Uint16Val.Val()
	assert.True(t, valid)
	assert.Equal(t, uint16(65535), uint16Val)

	uint32Val, valid := results[0].Uint32Val.Val()
	assert.True(t, valid)
	assert.Equal(t, uint32(4294967295), uint32Val)

	uint64Val, valid := results[0].Uint64Val.Val()
	assert.True(t, valid)
	assert.Equal(t, uint64(1234567890), uint64Val)

	float32Val, valid := results[0].Float32.Val()
	assert.True(t, valid)
	assert.InDelta(t, float32(3.14), float32Val, 0.01)

	float64Val, valid := results[0].Float64.Val()
	assert.True(t, valid)
	assert.InDelta(t, 2.718281828, float64Val, 0.00001)

	assert.NoError(t, mock.ExpectationsWereMet())
}
