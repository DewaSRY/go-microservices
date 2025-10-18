package types

import (
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string             `bson:"userID"`
	PackageSlug       string             `bson:"packageSlug"` // ex : van, luxury. sedan
	TotalPriceInCents float64            `bson:"totalPRiceInCents"`
	ExpiresAt         time.Time          `bson:"expiresAt"`
	Route             *types.Routes      `bson:"route"`
}

func (t *RideFareModel) ToTripProto() *tripgrpc.RideFare {
	return &tripgrpc.RideFare{
		Id:                t.ID.Hex(),
		UserID:            t.UserID,
		PackageSlug:       t.PackageSlug,
		TotalPriceInCents: t.TotalPriceInCents,
	}
}
