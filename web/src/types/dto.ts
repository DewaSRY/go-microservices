import { RouteFare } from "./types";

export interface RequestRideProps {
  pickup: [number, number];
  destination: [number, number];
}

export interface HTTPTripStartResponse {
  tripID: string;
}

export interface TripPreview {
  tripID: string;
  route: [number, number][];
  rideFares: RouteFare[];
  duration: number;
  distance: number;
}
