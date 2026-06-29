package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCartServiceWithMock(t *testing.T) (*CartService, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	return NewCartService(db.New(sqlDB), sqlDB), mock, sqlDB
}

func cartRow(cartID, userID uuid.UUID) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows([]string{"id", "user_id", "status", "created_at", "updated_at"}).
		AddRow(cartID, userID, "active", now, now)
}

func cartItemRow(itemID, cartID, productID uuid.UUID, qty int32) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity", "created_at"}).
		AddRow(itemID, cartID, productID, qty, time.Now())
}

func productRow(productID uuid.UUID, available bool) *sqlmock.Rows {
	now := time.Now()
	return sqlmock.NewRows([]string{
		"id", "restaurant_id", "name", "description", "price",
		"is_available", "kitchen_quantity", "created_at", "updated_at", "external_id", "image_url",
	}).AddRow(
		productID, uuid.New(), "Salad", "Fresh", "5.50",
		available, int32(10), now, now, int32(1), "https://example.com/salad.png",
	)
}

func TestCartService_AddItem_createsCartWhenMissing(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	productID := uuid.New()
	cartID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM products`).
		WithArgs(productID).
		WillReturnRows(productRow(productID, true))
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`INSERT INTO carts`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`FROM cart_items`).
		WithArgs(cartID, productID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`INSERT INTO cart_items`).
		WithArgs(sqlmock.AnyArg(), cartID, productID, int32(2)).
		WillReturnRows(cartItemRow(uuid.New(), cartID, productID, 2))
	mock.ExpectCommit()

	err := svc.AddItem(context.Background(), userID, dto.AddCartItemRequest{
		ProductID: productID,
		Quantity:  2,
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_AddItem_usesExistingCart(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	productID := uuid.New()
	cartID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM products`).
		WithArgs(productID).
		WillReturnRows(productRow(productID, true))
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`FROM cart_items`).
		WithArgs(cartID, productID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`INSERT INTO cart_items`).
		WithArgs(sqlmock.AnyArg(), cartID, productID, int32(1)).
		WillReturnRows(cartItemRow(uuid.New(), cartID, productID, 1))
	mock.ExpectCommit()

	err := svc.AddItem(context.Background(), userID, dto.AddCartItemRequest{
		ProductID: productID,
		Quantity:  1,
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_AddItem_productNotFound(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	productID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM products`).
		WithArgs(productID).
		WillReturnError(sql.ErrNoRows)

	err := svc.AddItem(context.Background(), userID, dto.AddCartItemRequest{
		ProductID: productID,
		Quantity:  1,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrProductNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_AddItem_productNotAvailable(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	productID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM products`).
		WithArgs(productID).
		WillReturnRows(productRow(productID, false))

	err := svc.AddItem(context.Background(), userID, dto.AddCartItemRequest{
		ProductID: productID,
		Quantity:  1,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrProductNotAvailable)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_AddItem_insufficientStock(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	productID := uuid.New()
	cartID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`FROM products`).
		WithArgs(productID).
		WillReturnRows(productRow(productID, true))
	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`FROM cart_items`).
		WithArgs(cartID, productID).
		WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(int32(8)))

	err := svc.AddItem(context.Background(), userID, dto.AddCartItemRequest{
		ProductID: productID,
		Quantity:  3,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInsufficientStock)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_GetCart_emptyWhenNoCart(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	resp, err := svc.GetCart(context.Background(), userID)
	require.NoError(t, err)
	assert.Empty(t, resp.CartID)
	assert.Empty(t, resp.Items)
	assert.InDelta(t, float64(0), resp.TotalAmount, 0.001)
	assert.Equal(t, int32(0), resp.ItemsCount)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_GetCart_withItems(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`FROM cart_items ci`).
		WithArgs(cartID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "price", "quantity"}).
			AddRow(productID, "Pasta", "10.00", int32(2)))

	resp, err := svc.GetCart(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, cartID.String(), resp.CartID)
	require.Len(t, resp.Items, 1)
	assert.Equal(t, productID.String(), resp.Items[0].ProductID)
	assert.InDelta(t, 10.0, resp.Items[0].Price, 0.001)
	assert.Equal(t, int32(2), resp.Items[0].Quantity)
	assert.InDelta(t, 20.0, resp.Items[0].LineTotal, 0.001)
	assert.InDelta(t, 20.0, resp.TotalAmount, 0.001)
	assert.Equal(t, int32(2), resp.ItemsCount)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_GetCart_invalidPrice(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()

	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectQuery(`FROM cart_items ci`).
		WithArgs(cartID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "price", "quantity"}).
			AddRow(uuid.New(), "Bad", "free", int32(1)))

	_, err := svc.GetCart(context.Background(), userID)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_ClearCart(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectExec(`DELETE FROM cart_items`).
		WithArgs(cartID).
		WillReturnResult(sqlmock.NewResult(0, 2))

	err := svc.ClearCart(context.Background(), userID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_ClearCart_cartNotFound(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	err := svc.ClearCart(context.Background(), userID)
	require.ErrorIs(t, err, ErrCartNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_RemoveCartItem(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()
	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectExec(`DELETE FROM cart_items`).
		WithArgs(cartID, productID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := svc.RemoveCartItem(context.Background(), userID, productID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_RemoveCartItem_itemNotFound(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	cartID := uuid.New()
	productID := uuid.New()
	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(cartRow(cartID, userID))
	mock.ExpectExec(`DELETE FROM cart_items`).
		WithArgs(cartID, productID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := svc.RemoveCartItem(context.Background(), userID, productID)
	require.ErrorIs(t, err, ErrCartItemNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartService_RemoveCartItem_cartNotFound(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	svc, mock, _ := newCartServiceWithMock(t)

	mock.ExpectQuery(`SELECT id, user_id, status, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	err := svc.RemoveCartItem(context.Background(), userID, uuid.New())
	require.ErrorIs(t, err, ErrCartNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}
