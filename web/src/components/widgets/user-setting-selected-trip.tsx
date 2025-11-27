import { Button } from "../ui/button";

import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import useRiderStore from "@/hooks/store/use-rider-store";
import { useSocketContext } from "@/components/provider/socket-provider";
import { RiderEvents } from "@/contracts/common";

export default function UserSettingSelectedTrip() {
  const { currentLocation, destination } = useRiderStore();
  const { sendMessage } = useSocketContext();

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
        <Button onClick={handleStartTrip}>Start Trip</Button>
      </div>
    </div>
  );
}
