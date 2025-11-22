import { Button } from "@components/ui/button";

import useRiderStore from "@/hooks/store/use-rider-store";
export default function RiderSelectDestination() {
  const { destination, setDestination } = useRiderStore();

  return (
    <div>
      <div>
        <h2>Select Destination To Go</h2>
      </div>
      <div>
        <p>Tab the map </p>
      </div>
      <div className="flex justify-between gap-2">
        <Button disabled={destination === undefined}>Create the trip</Button>
        <Button onClick={setDestination.bind(null, undefined)}>Cancel</Button>
      </div>
    </div>
  );
}
