import { Coordinate, PackageSlug, Entity } from "@/types/common";
import { DriverEvents } from "./common";

export type DriverWsRequest = DriverInitConnection;

interface DriverInitConnection {
  type: DriverEvents.DRIVER_INIT;
  data: {
    location: Coordinate;
    packageSlug: PackageSlug;
    entity: Entity;
  };
}
