import { useEffect } from "react";
import { WEBSOCKET_URL } from "../constants";
import { CarPackageSlug } from "../types/types";
import {
  ServerWsMessage,
  TripEvents,
  isValidWsMessage,
  isValidTripEvent,
  BackendEndpoints,
} from "../contracts";
import { useRouter } from "next/navigation";
import useDriverStore from "./useDriverStore";

interface useDriverConnectionProps {
  location: {
    latitude: number;
    longitude: number;
  };
  geohash: string;
  userID: string;
  packageSlug: CarPackageSlug;
}

export const useDriverStreamConnection = ({
  location,
  geohash,
  userID,
  packageSlug,
}: useDriverConnectionProps) => {
  const driverStore = useDriverStore();
  const router = useRouter();

  useEffect(() => {
    if (!userID) return;

    const websocket = new WebSocket(
      `${WEBSOCKET_URL}${BackendEndpoints.WS_DRIVERS}?userID=${userID}&packageSlug=${packageSlug}`
    );

    driverStore.setWs(websocket);

    websocket.onopen = () => {
      if (location) {
        // Send initial location
        websocket.send(
          JSON.stringify({
            type: TripEvents.DriverLocation,
            data: {
              location,
              geohash,
            },
          })
        );
      }
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data) as ServerWsMessage;

      if (!message || !isValidWsMessage(message)) {
        driverStore.setError(
          `Unknown message type "${message}", allowed types are: ${Object.values(
            TripEvents
          ).join(", ")}`
        );
        return;
      }

      switch (message.type) {
        case TripEvents.DriverTripRequest:
          const trip = message.data?.trip ?? message.data;
          driverStore.setRequestTrip(trip);
          break;
        case TripEvents.DriverRegister:
          driverStore.setDriver(message.data);
          break;
        case TripEvents.PaymentEventComplete:
          driverStore.resetTripStatus();
          router.push("?payment=success");
          break;
      }

      if (isValidTripEvent(message.type)) {
        driverStore.setTripStatus(message.type);
      } else {
        driverStore.setError(
          `Unknown message type "${
            message.type
          }", allowed types are: ${Object.values(TripEvents).join(", ")}`
        );
      }
    };

    websocket.onclose = () => {
      console.log("WebSocket closed");
    };

    websocket.onerror = (event) => {
      driverStore.setError("WebSocket error occurred");
      console.error("WebSocket error:", event);
    };

    return () => {
      console.log("Closing WebSocket");
      if (websocket.readyState === WebSocket.OPEN) {
        websocket.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID]);
};
