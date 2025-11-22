import { create } from "zustand";
import { Coordinate } from "@/types/common";

const initState = {
  destination: undefined as Coordinate | undefined,
  currentLocation: {
    latitude: -10.171781,
    longitude: 123.636975,
  } as Coordinate,
};

type Action = {
  setDestination: (_data: Coordinate | undefined) => void;
  setLocation: (_data: Coordinate) => void;
};

const useRiderStore = create<typeof initState & Action>((set) => {
  function setDestination(data: Coordinate | undefined) {
    set((res) => {
      return {
        ...res,
        destination: data,
      };
    });
  }

  function setLocation(data: Coordinate) {
    set((res) => {
      return {
        ...res,
        currentLocation: data,
      };
    });
  }

  return {
    ...initState,
    setDestination,
    setLocation,
  };
});

export default useRiderStore;
