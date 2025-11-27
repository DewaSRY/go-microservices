import { RiderFlowEvent } from "@/types/events";

import useRiderProfile from "@/hooks/state/useUserRideProfile";
import useRiderStore from "@/hooks/store/use-rider-store";
import { useEffect, useState } from "react";
import { Entity } from "@/types/common";

type eventState = RiderFlowEvent | undefined;

export default function useUserFlow() {
  const { mode } = useRiderProfile();
  const { destination, setDestination } = useRiderStore();

  const [currentEvent, setCurrentEvent] = useState<eventState>(undefined);
  const [isLockDestination, setLockDestination] = useState<boolean>(false);

  useEffect(() => {
    if (mode === undefined) {
      setCurrentEvent(undefined);
      setLockDestination(true);
      setDestination(undefined);
    }

    if (mode === Entity.RIDER) {
      setCurrentEvent(RiderFlowEvent.RIDER_INIT_CONNECTION);
      setLockDestination(false);
    }

    if (mode === Entity.RIDER && destination) {
      setCurrentEvent(RiderFlowEvent.TRIP_REQUESTED);
      setLockDestination(false);
    }

    if (mode === Entity.RIDER && destination && isLockDestination) {
      setCurrentEvent(RiderFlowEvent.WAITING_FOR_DRIVER);
      setLockDestination(true);
    }
  }, [mode, destination]);

  function setLockDestinationState(value: boolean) {
    if (mode !== Entity.RIDER || destination === undefined) return;
    setLockDestination(value);
  }

  return {
    currentEvent,
    setLockDestinationState,
    isLockDestination,
  };
}
