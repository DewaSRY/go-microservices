import { Coordinate, PackageSlug, Entity } from "@/types/common";
import { DriverEvents } from "./common";

export type DriverWsRequest = DriverInitConnection | DriverAcceptedTransaction;

interface DriverInitConnection {
  type: DriverEvents.DRIVER_INIT;
  data: {
    location: Coordinate;
    packageSlug: PackageSlug;
    entity: Entity;
  };
}

interface DriverAcceptedTransaction {
  type: DriverEvents.DRIVER_ACCEPTED_TRANSACTION;
  data: {
    transactionId: string;
  };
}
