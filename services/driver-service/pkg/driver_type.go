package pkg

import drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"

type ActiveDriver struct {
	Driver Driver
}

type Driver struct {
	ID             string
	Name           string
	ProfilePicture string
	CarPlage       string
	Geohash        string
	PackageSlug    string
	Location       *Location
}

func (t *Driver) ToDriverProto() *drivergrpc.Driver {
	return &drivergrpc.Driver{
		Id:             t.ID,
		Name:           t.Name,
		ProfilePicture: t.ProfilePicture,
		CarPlate:       t.CarPlage,
		Geohash:        t.Geohash,
		PackageSlug:    t.PackageSlug,
		Location:       t.Location.ToDriverProto(),
	}
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func (t *Location) ToDriverProto() *drivergrpc.Location {
	return &drivergrpc.Location{
		Latitude:  t.Latitude,
		Longitude: t.Longitude,
	}
}
