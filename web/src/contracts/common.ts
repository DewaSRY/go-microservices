export enum RiderEvents {
  CONNECTION_SUCCESS = "user.init-success.response",

  RIDER_CREATE_TRIP = "rider.create-trip.request",
  RIDER_CREATE_TRANSACTION = "rider.create-transaction.request",
}

export enum RideEvents {
  ROUTE_FOUND_RESPONSE = "trip.route-found.response",
  DRIVER_ACTIVE_RESPONSE = "driver.active-notify.response",
  RIDER_CREATE_TRANSACTION_RESPONSE = "rider.create-transaction.response",
  TRANSACTION_ACCEPTED_RESPONSE = "transaction.accepted.response",
}

export enum DriverEvents {
  DRIVER_INIT = "driver.init.request",
  DRIVER_ACCEPTED_TRANSACTION = "driver.transaction.accepted.request",
}
