package validation

import (
	"testing"

	bitedashv1 "bitedash/internal/pb/bitedash/v1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateGetOrderByIDRequest(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()
	tests := []struct {
		name string
		req  *bitedashv1.GetOrderByIDRequest
		code codes.Code
	}{
		{name: "valid", req: &bitedashv1.GetOrderByIDRequest{OrderId: orderID.String()}, code: codes.OK},
		{name: "nil request", req: nil, code: codes.InvalidArgument},
		{name: "empty id", req: &bitedashv1.GetOrderByIDRequest{}, code: codes.InvalidArgument},
		{name: "invalid id", req: &bitedashv1.GetOrderByIDRequest{OrderId: "invalid"}, code: codes.InvalidArgument},
		{name: "nil uuid", req: &bitedashv1.GetOrderByIDRequest{OrderId: uuid.Nil.String()}, code: codes.InvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ValidateGetOrderByIDRequest(tt.req)
			assert.Equal(t, tt.code, status.Code(err))
			if tt.code == codes.OK {
				assert.Equal(t, orderID, actual)
			}
		})
	}
}

func TestValidateAddCartItemRequest(t *testing.T) {
	t.Parallel()

	productID := uuid.New()
	tests := []struct {
		name string
		req  *bitedashv1.AddCartItemRequest
		code codes.Code
	}{
		{name: "valid", req: &bitedashv1.AddCartItemRequest{ProductId: productID.String(), Quantity: 2}, code: codes.OK},
		{name: "nil request", req: nil, code: codes.InvalidArgument},
		{name: "empty product", req: &bitedashv1.AddCartItemRequest{Quantity: 1}, code: codes.InvalidArgument},
		{name: "invalid product", req: &bitedashv1.AddCartItemRequest{ProductId: "invalid", Quantity: 1}, code: codes.InvalidArgument},
		{name: "zero quantity", req: &bitedashv1.AddCartItemRequest{ProductId: productID.String()}, code: codes.InvalidArgument},
		{name: "negative quantity", req: &bitedashv1.AddCartItemRequest{ProductId: productID.String(), Quantity: -1}, code: codes.InvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := ValidateAddCartItemRequest(tt.req)
			assert.Equal(t, tt.code, status.Code(err))
			if tt.code == codes.OK {
				assert.Equal(t, productID, actual.ProductID)
				assert.Equal(t, int32(2), actual.Quantity)
			}
		})
	}
}

func TestValidateListRestaurantsRequest(t *testing.T) {
	t.Parallel()

	require.NoError(t, ValidateListRestaurantsRequest(&bitedashv1.ListRestaurantsRequest{}))
	require.NoError(t, ValidateListRestaurantsRequest(&bitedashv1.ListRestaurantsRequest{Page: 1, Limit: 50}))
	assert.Equal(t, codes.InvalidArgument, status.Code(ValidateListRestaurantsRequest(nil)))
	assert.Equal(t, codes.InvalidArgument, status.Code(ValidateListRestaurantsRequest(
		&bitedashv1.ListRestaurantsRequest{Page: -1},
	)))
	assert.Equal(t, codes.InvalidArgument, status.Code(ValidateListRestaurantsRequest(
		&bitedashv1.ListRestaurantsRequest{Limit: -1},
	)))
}
