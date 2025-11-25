import { Button } from "@components/ui/button";
import useRiderStore from "@/hooks/store/use-rider-store";
import { useSocketContext } from "@components/provider/socket-provider";
import { RideEvents, RiderEvents } from "@/contracts/common";

export default function RiderSelectDestination() {
  const { sendMessage } = useSocketContext();
  const { destination, setDestination, currentLocation } = useRiderStore();

  function handleCreateTrip() {
    if (destination === undefined) return;
    sendMessage({
      type: RiderEvents.TRIP_CREATE_EVENT,
      data: {
        destination: destination,
        pickup: currentLocation,
      },
    });
  }

  return (
    <div>
      <div>
        <h2>Select Destination To Go</h2>
      </div>
      <div>
        <p>Tab the map </p>
      </div>
      <div className="flex justify-between gap-2">
        <Button disabled={destination === undefined} onClick={handleCreateTrip}>
          Create the trip
        </Button>
        <Button onClick={setDestination.bind(null, undefined)}>Cancel</Button>
      </div>
    </div>
  );
}
