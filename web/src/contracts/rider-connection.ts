import { Coordinate, PackageSlug, Entity } from "@/types/common";

export enum RiderEvents {
  InitConnection = "user.event.init",
}

//Messages sent from client to server via web socket
export type RiderWsMessage = RiderInitConnection;

interface RiderInitConnection {
  type: RiderEvents.InitConnection;
  data: {
    location: Coordinate;
    packageSlug: PackageSlug;
    entity: Entity;
  };
}
