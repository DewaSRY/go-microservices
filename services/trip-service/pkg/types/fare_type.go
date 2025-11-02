package types

// import (
// 	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
// 	"DewaSRY/go-microservices/shared/types"
// 	"fmt"
// 	"time"
// )

// type FareType struct {
// 	ID                uint64
// 	UserID            uint64
// 	PackageSlug       string
// 	TotalPriceInCents float64
// 	ExpiresAt         time.Time
// 	Route             *types.Routes
// }

// func (t *FareType) ToTripProto() *tripgrpc.RideFare {
// 	return &tripgrpc.RideFare{
// 		Id:                fmt.Sprintf("{%d}", t.ID),
// 		UserID:            fmt.Sprintf("{%d}", t.UserID),
// 		PackageSlug:       t.PackageSlug,
// 		TotalPriceInCents: t.TotalPriceInCents,
// 	}
// }
