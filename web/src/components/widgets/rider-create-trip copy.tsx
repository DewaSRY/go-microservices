import { Button } from "../ui/button";

import useRiderStore from "@/hooks/store/use-rider-store";
import { useSocketContext } from "@/components/provider/socket-provider";
import { RiderEvents } from "@/contracts/common";
import { useFlowContext } from "@components/provider/user-flow-provider";
import { RiderFlowEvent } from "@/types/events";

export default function RiderCreateTrip() {
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

  return (
    <div>
      <div>Selected Trip Settings Component</div>
      <div>
        <Button disabled={destination === undefined} onClick={handleStartTrip}>
          Start Trip
        </Button>
      </div>
    </div>
  );
}
