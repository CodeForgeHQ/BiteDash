package service

import (
	"context"
	"database/sql"
	"testing"

	db "bitedash/internal/db/sqlc"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newOrderServiceWithMock(t *testing.T) (*OrderService, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	return NewOrderService(sqlDB, db.New(sqlDB)), mock
}

func cartItemsWithStockRows(productID uuid.UUID) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"product_id", "name", "price", "quantity", "is_available", "kitchen_quantity"}).
		AddRow(productID, "Salad", "5.50", int32(2), true, int32(10))
}

func TestOrderService_MakeOrder_success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()

	svc, mock := newOrderServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`p.kitchen_quantity`).
		WithArgs(cartID).
		WillReturnRows(cartItemsWithStockRows(productID))
	mock.ExpectExec(`UPDATE products`).
		WithArgs(productID, int32(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO orders`).
		WithArgs(sqlmock.AnyArg(), userID, "pending", "11.00").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO order_items`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), productID, "Salad", "5.50", int32(2), "11.00").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM cart_items`).
		WithArgs(cartID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	resp, err := svc.MakeOrder(context.Background(), userID)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.OrderID)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, 11.0, resp.TotalAmount)
	assert.Equal(t, int32(2), resp.ItemsCount)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderService_MakeOrder_cartNotFound(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	svc, mock := newOrderServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	_, err := svc.MakeOrder(context.Background(), userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCartNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderService_MakeOrder_emptyCart(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()

	svc, mock := newOrderServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`p.kitchen_quantity`).
		WithArgs(cartID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "price", "quantity", "is_available", "kitchen_quantity"}))

	_, err := svc.MakeOrder(context.Background(), userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCartEmpty)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderService_MakeOrder_invalidPrice(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()

	svc, mock := newOrderServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`p.kitchen_quantity`).
		WithArgs(cartID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "price", "quantity", "is_available", "kitchen_quantity"}).
			AddRow(uuid.New(), "Mystery", "N/A", int32(1), true, int32(10)))

	_, err := svc.MakeOrder(context.Background(), userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid price")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderService_MakeOrder_insufficientStock(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()

	svc, mock := newOrderServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`p.kitchen_quantity`).
		WithArgs(cartID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "price", "quantity", "is_available", "kitchen_quantity"}).
			AddRow(productID, "Salad", "5.50", int32(5), true, int32(2)))

	_, err := svc.MakeOrder(context.Background(), userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInsufficientStock)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderService_MakeOrder_productNotAvailable(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()

	svc, mock := newOrderServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`p.kitchen_quantity`).
		WithArgs(cartID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "price", "quantity", "is_available", "kitchen_quantity"}).
			AddRow(productID, "Salad", "5.50", int32(1), false, int32(10)))

	_, err := svc.MakeOrder(context.Background(), userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrProductNotAvailable)
	require.NoError(t, mock.ExpectationsWereMet())
}
