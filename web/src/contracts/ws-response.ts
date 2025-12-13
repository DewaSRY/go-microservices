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
  | RiderCreateTransaction
  | TransactionAcceptedResponse;

interface RiderInitConnectionSuccess {
  type: RiderEvents.CONNECTION_SUCCESS;
  data: {
    message: string;
  };
}

interface TripFound {
  type: RideEvents.ROUTE_FOUND_RESPONSE;
  data: RouteData;
}

interface DriverActiveList {
  type: RideEvents.DRIVER_ACTIVE_RESPONSE;
  data: DriverActiveRecord[];
}

interface RiderCreateTransaction {
  type: RideEvents.RIDER_CREATE_TRANSACTION_RESPONSE;
  data: RiderCreateTransactionResponse;
}

interface TransactionAcceptedResponse {
  type: RideEvents.TRANSACTION_ACCEPTED_RESPONSE;
  data: {
    transactionId: string;
  };
}
