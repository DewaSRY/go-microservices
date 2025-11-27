/**
 TODO: think about negative event
*/

export enum RiderFlowEvent {
  /**
   * After connection mark as rider it will
   */
  RIDER_INIT_CONNECTION = "RIDER_INIT_CONNECTION",
  TRIP_REQUESTED = "RIDER_TRIP_REQUESTED",
  WAITING_FOR_DRIVER = "WAITING_FOR_DRIVER",
}
