import { Button } from "../ui/button";

import useRiderStore from "@/hooks/store/use-rider-store";
import { useSocketContext } from "@/components/provider/socket-provider";
import { useFlowContext } from "@components/provider/user-flow-provider";
import { RiderFlowEvent } from "@/types/events";
import DriverActiveCard from "./driver-active-card";

export default function RiderListingDriver() {
  const { currentEvent } = useFlowContext();
  const { currentLocation, destination, setDestination } = useRiderStore();
  const { resetRoute, driverActive } = useSocketContext();

  function handleCancelTrip() {
    if (currentLocation && destination) {
      setDestination(undefined);
      resetRoute();
    }
  }

  return (
    <div>
      <div>this is supper</div>

      <div>
        {driverActive.length === 0 ? (
          <div>
            <span>waiting for driver</span>
          </div>
        ) : (
          <div>
            {driverActive.map((x, idx) => (
              <DriverActiveCard
                key={idx}
                title={`car-${idx}`}
                packageSlug={x.packageSlug}
                driverId={x.driverId}
              />
            ))}
          </div>
        )}
      </div>
      <div>
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
