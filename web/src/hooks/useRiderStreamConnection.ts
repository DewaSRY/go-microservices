import { useEffect } from "react";
import { WEBSOCKET_URL } from "../constants";
import { Coordinate } from "../types/types";
import {
  TripEvents,
  ServerWsMessage,
  isValidWsMessage,
  BackendEndpoints,
} from "../contracts";
import { useRouter } from "next/navigation";
import useRiderStore from "./useRiderStore";

export function useRiderStreamConnection(location: Coordinate, userID: string) {
  const riderStore = useRiderStore();
  const router = useRouter();

  useEffect(() => {
    if (!userID) return;

    const ws = new WebSocket(
      `${WEBSOCKET_URL}${BackendEndpoints.WS_RIDERS}?userID=${userID}`
    );
    riderStore.setWs(ws);

    ws.onopen = () => {
      // Send initial location
      if (location) {
        ws.send(
          JSON.stringify({
            type: TripEvents.DriverLocation,
            data: {
              location,
            },
          })
        );
      }
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data) as ServerWsMessage;

      if (!message || !isValidWsMessage(message)) {
        riderStore.setError(
          `Unknown message type "${message}", allowed types are: ${Object.values(
            TripEvents
          ).join(", ")}`
        );
        return;
      }

      switch (message.type) {
        case TripEvents.DriverLocation:
          riderStore.setDrivers(message.data);
          break;
        case TripEvents.PaymentSessionCreated:
          riderStore.setPaymentSession(message.data);
          riderStore.setTripStatus(message.type);
          break;
        case TripEvents.DriverAssigned:
          riderStore.setAssignedDriver(message.data.driver ?? null);
          riderStore.setTripStatus(message.type);
          break;
        case TripEvents.Created:
          riderStore.setTripStatus(message.type);
          break;
        case TripEvents.NoDriversFound:
          riderStore.setTripStatus(message.type);
          break;
        case TripEvents.PaymentEventComplete:
          router.push("?payment=success");
          break;
      }
    };

    ws.onclose = () => {
      console.log("WebSocket closed");
    };

    ws.onerror = (event) => {
      riderStore.setError("WebSocket error occurred");
      console.error("WebSocket error:", event);
    };

    return () => {
      console.log("Closing WebSocket");
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID]);
}
