import { create } from "zustand";
import { Coordinate } from "@/types/common";

const USER_LOCATION = "USER_LOCATION";

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
  reset: () => void;

  initData: () => void;
};

const useRideStore = create<typeof initState & Action>((set) => {
  function setDestination(data: Coordinate | undefined) {
    set((res) => {
      return {
        ...res,
        destination: data,
      };
    });
  }

  function setLocation(data: Coordinate) {
    localStorage.setItem(USER_LOCATION, JSON.stringify(data));

    set((res) => {
      return {
        ...res,
        currentLocation: data,
      };
    });
  }

  function reset() {
    set((prev) => {
      prev.destination = undefined;
      return prev;
    });
  }

  function initData() {
    let userLocation: Coordinate = initState.currentLocation;

    const storedLocation = localStorage.getItem(USER_LOCATION);

    if (storedLocation) {
      userLocation = JSON.parse(storedLocation);
    }

    set((res) => {
      return {
        ...res,
        currentLocation: initState.currentLocation,
      };
    });
  }

  return {
    ...initState,
    setDestination,
    setLocation,
    reset,
    initData,
  };
});

export default useRideStore;
