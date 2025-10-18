package triptype

import (
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Update this letter
type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare RideFareModel
	Driver   *TripDriver
}

func (t *TripModel) ToTripProto() *tripgrpc.Trip {
	return &tripgrpc.Trip{
		Id:           t.ID.Hex(),
		UserID:       t.UserID,
		SelectedFare: t.RideFare.ToTripProto(),
		Status:       t.Status,
		Route:        t.RideFare.Route.ToRouteProto(),
	}
}
