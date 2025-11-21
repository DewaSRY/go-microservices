import { Coordinate, PackageSlug, Entity } from "@/types/common";
import { RiderEvents } from "./common";

//Messages sent from client to server via web socket
export type RiderWsRequest = RiderInitConnection;

interface RiderInitConnection {
  type: RiderEvents.INIT_CONNECTION;
  data: {
    location: Coordinate;
    packageSlug: PackageSlug;
    entity: Entity;
  };
}
