package types

import tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"

type TripDriver struct {
	ID             string
	Name           string
	ProfilePicture string
	CarPlate       string
}

func (t *TripDriver) ToTripProto() *tripgrpc.TripDriver {
	return &tripgrpc.TripDriver{
		Id:             t.ID,
		Name:           t.Name,
		ProfilePicture: t.ProfilePicture,
	}
}
