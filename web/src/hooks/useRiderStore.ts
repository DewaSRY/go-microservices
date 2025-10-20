import {
  ClientWsMessage,
  PaymentEventSessionCreatedData,
  TripEvents,
} from "@/contracts";
import { Driver, Trip } from "@/types/types";
import { create } from "zustand";

//  const [tripStatus, setTripStatus] = useState<TripEvents | null>(null);
const initState = {
  drivers: [] as Driver[],
  tripStatus: null as TripEvents | null,
  paymentSession: null as PaymentEventSessionCreatedData | null,
  assignedDriver: null as Driver | null,
  ws: null as WebSocket | null,
  driver: null as Driver | null,
  error: null as string | null,
};

type Action = {
  setError: (error: string | null) => void;
  setWs: (ws: WebSocket | null) => void;
  setPaymentSession: (data: PaymentEventSessionCreatedData | null) => void;
  setAssignedDriver: (data: Driver | null) => void;
  setDrivers: (data: Driver[]) => void;
  setTripStatus: (data: TripEvents | null) => void;
  sendMessage: (message: ClientWsMessage) => void;
  resetTripStatus: () => void;
};

const useRiderStore = create<typeof initState & Action>((set, get) => {
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

  const setPaymentSession = (data: PaymentEventSessionCreatedData | null) => {
    set((state) => ({
      ...state,
      paymentSession: data,
    }));
  };

  const setAssignedDriver = (data: Driver | null) => {
    set((state) => ({
      ...state,
      driver: data,
    }));
  };

  const setDrivers = (data: Driver[]) => {
    set((state) => ({
      ...state,
      drivers: data,
    }));
  };

  const setTripStatus = (data: TripEvents | null) => {
    set((state) => ({
      ...state,
      tripStatus: data,
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
      tripStatus: null,
    }));
  };

  return {
    ...initState,
    setError,
    setWs,
    setPaymentSession,
    setAssignedDriver,
    setDrivers,
    setTripStatus,
    sendMessage,
    resetTripStatus,
  };
});

export default useRiderStore;
