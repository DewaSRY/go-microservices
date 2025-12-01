import { Button } from "../ui/button";

import useRiderStore from "@/hooks/store/use-rider-store";
import { useSocketContext } from "@/components/provider/socket-provider";
import { RiderEvents } from "@/contracts/common";
import { useFlowContext } from "@components/provider/user-flow-provider";
import { RiderFlowEvent } from "@/types/events";

export default function UserSettingSelectedTrip() {
  const { currentEvent } = useFlowContext();
  const { currentLocation, destination, setDestination } = useRiderStore();
  const { sendMessage } = useSocketContext();
  const { resetRoute } = useSocketContext();

  function handleStartTrip() {
    if (currentLocation && destination) {
      sendMessage({
        type: RiderEvents.RIDER_CREATE_TRIP,
        data: {
          pickup: currentLocation,
          destination: destination,
        },
      });
    }
  }

  function handleCancelTrip() {
    if (currentLocation && destination) {
      setDestination(undefined);
      resetRoute();
    }
  }

  return (
    <div>
      <div>Selected Trip Settings Component</div>
      <div>
        <Button
          disabled={currentEvent !== RiderFlowEvent.TRIP_REQUESTED}
          onClick={handleStartTrip}
        >
          Start Trip
        </Button>

        <Button
          disabled={currentEvent !== RiderFlowEvent.WAITING_FOR_DRIVER}
          onClick={handleCancelTrip}
        >
          Cancel Trip
        </Button>
      </div>
    </div>
  );
}
