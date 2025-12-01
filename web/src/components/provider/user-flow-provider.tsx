import { RiderWsRequest } from "@/contracts/ws-request";
import { RiderWsResponse } from "@/contracts/ws-response";
import { useCallback, useMemo, useRef, useState } from "react";
import { useEffect } from "react";
import {
  PropsWithChildren,
  useContext,
  HTMLAttributes,
  createContext,
} from "react";
import { Entity } from "@/types/common";
import { RiderFlowEvent } from "@/types/events";

import useRiderProfile from "@/hooks/state/useUserRideProfile";
import useRiderStore from "@/hooks/store/use-rider-store";
type eventState = RiderFlowEvent | undefined;
const UserFlowProvider = createContext({
  currentEvent: undefined as eventState | undefined,
  setLockDestinationState: (val: boolean) => {},
  isLockDestination: false,
  setIsHaveRideRoute: (val: boolean) => {},
});

UserFlowProvider.displayName = "user-flow";

interface ProviderPops
  extends HTMLAttributes<HTMLDivElement>,
    PropsWithChildren {}

export default function Provider({ children }: ProviderPops) {
  const { mode } = useRiderProfile();
  const { destination, setDestination } = useRiderStore();

  const [isLockDestination, setLockDestination] = useState<boolean>(false);
  const [isHaveRideRoute, setIsHaveRideRoute] = useState<boolean>(false);

  const currentEvent = useMemo(() => {
    if (mode !== Entity.RIDER) return undefined;
    console.log("is have route", isHaveRideRoute);
    if (isHaveRideRoute) {
      return RiderFlowEvent.WAITING_FOR_DRIVER;
    }
    if (isHaveRideRoute === false) {
      return RiderFlowEvent.TRIP_REQUESTED;
    }
    return undefined;
  }, [mode, isHaveRideRoute]);

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

  useEffect(() => {
    console.log("current event is", currentEvent);
  });

  function _setIsHaveRideRoute(val: boolean) {
    setIsHaveRideRoute(val);
  }

  function setLockDestinationState(value: boolean) {
    if (mode !== Entity.RIDER || destination === undefined) return;
    setLockDestination(value);
  }
  return (
    <UserFlowProvider.Provider
      value={{
        currentEvent,
        isLockDestination,
        setIsHaveRideRoute: _setIsHaveRideRoute,
        setLockDestinationState,
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
