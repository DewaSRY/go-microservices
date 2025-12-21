export type Coordinate = {
  latitude: number;
  longitude: number;
};

export enum PackageSlug {
  SEDAN = "sedan",
  SUV = "suv",
  VAN = "van",
  LUXURY = "luxury",
}

export enum Entity {
  RIDER = "RIDER",
  DRIVER = "DRIVER",
}

export type RouteData = {
  coordinate: Coordinate[];
  distance: number;
  duration: number;
};

export type DriverActiveRecord = {
  coordinate: Coordinate;
  packageSlug: PackageSlug;
  driverId: string;
};

export type RiderCreateTransactionResponse = {
  transactionId: string;
};

export type LocationPermissionState =
  | "granted"
  | "prompt"
  | "denied"
  | "unsupported";
