import { Coordinate, PackageSlug, Entity } from "@/types/common";
import { RiderEvents } from "./common";

export type RiderWsRequest = RiderCreateTripRequest;

// interface RiderInitConnection {
//   type: RiderEvents.INIT_CONNECTION;
//   data: {
//     location: Coordinate;
//     packageSlug: PackageSlug;
//     entity: Entity;
//   };
// }

// interface TripCreateRequest {
//   type: RiderEvents.TRIP_CREATE_EVENT;
//   data: {
//     pickup: Coordinate;
//     destination: Coordinate;
//   };
// }

interface RiderCreateTripRequest {
  type: RiderEvents.RIDER_CREATE_TRIP;
  data: {
    pickup: Coordinate;
    destination: Coordinate;
  };
}
