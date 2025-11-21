import { RiderEvents } from "./common";

//Messages sent from client to server via web socket
export type RiderWsResponse = RiderInitConnectionSuccess;

interface RiderInitConnectionSuccess {
  type: RiderEvents.CONNECTION_SUCCESS;
  data: {
    message: string;
  };
}
