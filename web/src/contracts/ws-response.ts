import {
  DriverActiveRecord,
  RouteData,
  RiderCreateTransactionResponse,
} from "@/types/common";
import { RiderEvents, RideEvents } from "./common";

//Messages sent from client to server via web socket
export type RiderWsResponse =
  | RiderInitConnectionSuccess
  | TripFound
  | DriverActiveList
  | RiderCreateTransaction;

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

interface RiderCreateTransaction {
  type: RideEvents.RIDER_CREATE_TRANSACTION;
  data: RiderCreateTransactionResponse;
}
