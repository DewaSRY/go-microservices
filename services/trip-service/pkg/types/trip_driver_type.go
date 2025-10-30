package types

import (
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"fmt"
)

type TripDriverType struct {
	ID             uint64
	Name           string
	ProfilePicture string
	CarPlate       string
}

func (t *TripDriverType) ToTripProto() *tripgrpc.TripDriver {
	return &tripgrpc.TripDriver{
		Id:             fmt.Sprintf("{%d}", t.ID),
		Name:           t.Name,
		ProfilePicture: t.ProfilePicture,
	}
}
