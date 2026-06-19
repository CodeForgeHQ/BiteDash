package mapper

import (
	db "bitedash/internal/db/sqlc"
	bitedashv1 "bitedash/internal/pb/bitedash/v1"
)

func UserToProto(user db.User) *bitedashv1.User {
	return &bitedashv1.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}
}
