export enum RiderEvents {
  CONNECTION_SUCCESS = "user.init-success.response",

  RIDER_CREATE_TRIP = "rider.create-trip.request",
}

export enum RideEvents {
  ROUTE_FOUND = "trip.route-found.response",
  DRIVER_ACTIVE = "driver.active-notify.response",
}

export enum DriverEvents {
  DRIVER_INIT = "driver.init.request",
}
