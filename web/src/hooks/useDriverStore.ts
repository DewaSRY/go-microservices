import { create } from "zustand";
import { Driver, Trip } from "@/types/types";
import { ClientWsMessage, TripEvents } from "@/contracts";

const initState = {
  trip: null as Trip | null,
  tripEvent: null as TripEvents | null,
  error: null as string | null,
  ws: null as WebSocket | null,
  driver: null as Driver | null,
};

type Action = {
  setRequestTrip: (trip: Trip | null) => void;
  setTripStatus: (event: TripEvents | null) => void;
  setError: (error: string | null) => void;
  setWs: (ws: WebSocket | null) => void;
  setDriver: (driver: Driver | null) => void;
  sendMessage: (message: ClientWsMessage) => void;
  resetTripStatus: () => void;
};

const useDriverStore = create<typeof initState & Action>((set, get) => {
  const setRequestTrip = (trip: Trip | null) => {
    set((state) => ({
      ...state,
      trip: trip,
    }));
  };

  const setTripStatus = (event: TripEvents | null) => {
    set((state) => ({
      ...state,
      tripEvent: event,
    }));
  };

  const setError = (error: string | null) => {
    set((state) => ({
      ...state,
      error: error,
    }));
  };

  const setWs = (ws: WebSocket | null) => {
    set((state) => ({
      ...state,
      ws: ws,
    }));
  };

  const setDriver = (driver: Driver | null) => {
    set((state) => ({
      ...state,
      driver: driver,
    }));
  };

  const sendMessage = (message: ClientWsMessage) => {
    const { ws, setError } = get();
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    } else {
      setError("WebSocket is not connected");
    }
  };

  const resetTripStatus = () => {
    set((state) => ({
      ...state,
      tripEvent: null,
      trip: null,
      driver: null,
      error: null,
    }));
  };

  return {
    ...initState,
    setRequestTrip,
    setTripStatus,
    setError,
    setWs,
    setDriver,
    sendMessage,
    resetTripStatus,
  };
});

export default useDriverStore;
