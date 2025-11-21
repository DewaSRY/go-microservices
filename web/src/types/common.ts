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
