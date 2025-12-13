import { Coordinate, PackageSlug, Entity } from "@/types/common";
import { RiderEvents } from "./common";

export type RiderWsRequest =
  | RiderCreateTripRequest
  | RiderCreateTransactionRequest;

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

interface RiderCreateTransactionRequest {
  type: RiderEvents.RIDER_CREATE_TRANSACTION;
  data: {
    driverId: string;
  };
}
