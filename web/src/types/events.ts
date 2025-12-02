/**
 TODO: think about negative event
*/

export enum RiderFlowEvent {
  /**
   * After connection mark as rider it will
   */
  TRIP_REQUESTED = "RIDER_TRIP_REQUESTED",
  WAITING_FOR_DRIVER = "WAITING_FOR_DRIVER",
}

export enum DriverFlowEvents {
  /**
   * After connection mark as rider it will
   */
  DRIVER_INIT_CONN = "DRIVER_INIT_CONN",
  DRIVER_WAITING_FOR_RIDER = "DRIVER_WAITING_FOR_RIDER",
}
