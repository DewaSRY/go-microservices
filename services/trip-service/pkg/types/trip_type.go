package types

// import (
// 	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
// 	"fmt"
// )

// type TripType struct {
// 	ID       uint64
// 	UserID   uint64
// 	Status   string
// 	RideFare FareType
// 	Driver   *TripDriverType
// }

// func (t *TripType) ToTripProto() *tripgrpc.Trip {
// 	return &tripgrpc.Trip{
// 		Id:           fmt.Sprintf("{%d}", t.ID),
// 		UserID:       fmt.Sprintf("{%d}", t.UserID),
// 		SelectedFare: t.RideFare.ToTripProto(),
// 		Status:       t.Status,
// 		Route:        t.RideFare.Route.ToRouteProto(),
// 	}
// }
