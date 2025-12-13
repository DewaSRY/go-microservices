import type { RiderWsRequest } from "./ws-rider-request";
import type { DriverWsRequest } from "./ws-driver-request";

//Messages sent from client to server via web socket
export type UserWsRequest = RiderWsRequest | DriverWsRequest;
