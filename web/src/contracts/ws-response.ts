import { DriverActiveRecord, RouteData } from "@/types/common";
import { RiderEvents, RideEvents } from "./common";

//Messages sent from client to server via web socket
export type RiderWsResponse =
  | RiderInitConnectionSuccess
  | TripFound
  | DriverActiveList;

interface RiderInitConnectionSuccess {
  type: RiderEvents.CONNECTION_SUCCESS;
  data: {
    message: string;
  };
}

interface TripFound {
  type: RideEvents.ROUTE_FOUND;
  data: RouteData;
}

interface DriverActiveList {
  type: RideEvents.DRIVER_ACTIVE;
  data: DriverActiveRecord[];
}
