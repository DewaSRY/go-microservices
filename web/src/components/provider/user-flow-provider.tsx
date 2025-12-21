///
import { useCallback, useMemo, useRef, useState } from "react";
import { useEffect } from "react";
import {
  PropsWithChildren,
  useContext,
  HTMLAttributes,
  createContext,
} from "react";
import { Entity } from "@/types/common";
import {
  RiderFlowEvent,
  DriverFlowEvents,
  RideFlowEvents,
} from "@/types/events";

import useRiderProfile from "@/hooks/state/useUserRideProfile";
import useRiderStore from "@/hooks/store/use-ride-store";
import { useSocketContext } from "./socket-provider";
import { useUserLocationContext } from "./user-location-provider";

type eventState =
  | RideFlowEvents
  | RiderFlowEvent
  | DriverFlowEvents
  | undefined;

const UserFlowProvider = createContext({
  currentEvent: undefined as eventState | undefined,
  setLockDestinationState: (val: boolean) => {},
  isLockDestination: false,
  setIsHaveRideRoute: (val: boolean) => {},
  resetState: () => {},
});

UserFlowProvider.displayName = "user-flow";

interface ProviderPops
  extends HTMLAttributes<HTMLDivElement>,
    PropsWithChildren {}

export default function Provider({ children }: ProviderPops) {
  const { locationPermission } = useUserLocationContext();
  const { mode, packageSlug } = useRiderProfile();
  const { destination, setDestination } = useRiderStore();
  const { transactionId, isTransactionAccepted } = useSocketContext();

  const [isLockDestination, setLockDestination] = useState<boolean>(false);
  const [isHaveRideRoute, setIsHaveRideRoute] = useState<boolean>(false);

  const currentEvent = useMemo(() => {
    if (locationPermission !== "granted") {
      return RideFlowEvents.LOCATION_NOT_SET;
    }

    if (mode === undefined) {
      return RideFlowEvents.MODE_NOT_SET;
    }

    if (mode === Entity.RIDER) {
      if (isTransactionAccepted) {
        return RiderFlowEvent.RIDER_TRANSACTION_SUCCESS;
      }

      if (isHaveRideRoute && transactionId !== undefined) {
        return RiderFlowEvent.RIDER_WAITING_DRIVER_CONFIRMATION;
      }

      if (isHaveRideRoute) {
        return RiderFlowEvent.WAITING_FOR_DRIVER;
      }

      if (isHaveRideRoute === false) {
        return RiderFlowEvent.TRIP_REQUESTED;
      }
    }

    if (mode === Entity.DRIVER) {
      if (isTransactionAccepted) {
        return DriverFlowEvents.DRIVER_TRANSACTION_SUCCESS;
      }

      if (packageSlug !== undefined && transactionId !== undefined) {
        return DriverFlowEvents.RIDER_CREATE_TRANSACTION;
      }

      if (packageSlug !== undefined) {
        return DriverFlowEvents.DRIVER_WAITING_FOR_RIDER;
      }

      return DriverFlowEvents.DRIVER_INIT_CONN;
    }

    return undefined;
  }, [
    mode,
    isHaveRideRoute,
    packageSlug,
    transactionId,
    isTransactionAccepted,
    locationPermission,
  ]);

  useEffect(() => {
    if (mode === Entity.RIDER) {
      if (isHaveRideRoute === true) {
        setLockDestination(true);
        return;
      }

      if (isHaveRideRoute === false) {
        setLockDestination(false);
        return;
      }
    }

    if (mode === undefined) {
      setLockDestination(true);
      setDestination(undefined);
    }
  }, [mode, destination, isHaveRideRoute]);

  function _setIsHaveRideRoute(val: boolean) {
    setIsHaveRideRoute(val);
  }

  function setLockDestinationState(value: boolean) {
    if (mode !== Entity.RIDER || destination === undefined) return;
    setLockDestination(value);
  }

  function resetState() {
    setLockDestination(true);
    setIsHaveRideRoute(false);
  }

  return (
    <UserFlowProvider.Provider
      value={{
        currentEvent,
        isLockDestination,
        setIsHaveRideRoute: _setIsHaveRideRoute,
        setLockDestinationState,
        resetState,
      }}
    >
      {children}
    </UserFlowProvider.Provider>
  );
}

export function useFlowContext() {
  const context = useContext(UserFlowProvider);

  if (!context) throw Error("hook_use_outside_user_flow_the_context");

  return context;
}
