export enum RiderEvents {
  INIT_CONNECTION = "user.init.request",
  CONNECTION_SUCCESS = "user.init-success.response",

  TRIP_CREATE_EVENT = "trip_flow.create.request",

  RIDER_CREATE_TRIP = "rider.create-trip.request",
}

export enum RideEvents {
  ROUTE_FOUND = "trip.route-found.response",
}
