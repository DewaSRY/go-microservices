export enum RiderEvents {
  CONNECTION_SUCCESS = "user.init-success.response",

  RIDER_CREATE_TRIP = "rider.create-trip.request",

  RIDER_CREATE_TRANSACTION = "rider.create-transaction.request",
}

export enum RideEvents {
  ROUTE_FOUND = "trip.route-found.response",
  DRIVER_ACTIVE = "driver.active-notify.response",
  RIDER_CREATE_TRANSACTION = "rider.create-transaction.response",
}

export enum DriverEvents {
  DRIVER_INIT = "driver.init.request",
}
