import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import useRiderStore from "@/hooks/store/use-ride-store";
import { useSocketContext } from "@/components/provider/socket-provider";
import { useFlowContext } from "@components/provider/user-flow-provider";
import { RiderFlowEvent } from "@/types/events";
import DriverActiveCard from "./driver-active";

import TitleBadge from "@/components/ui/title-badge";
import { Info } from "lucide-react";

export default function RiderListingDriver() {
  const translate = useTranslations("ride.listingDriver");
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
      <TitleBadge>{translate("title")}</TitleBadge>

      <br className="h-16" />

      <div>
        {driverActive.length === 0 ? (
          <div>
            <div className="w-11/12 mx-auto rounded-xl bg-yellow-400/10 px-4 py-2 text-gray-700 flex items-center gap-2">
              <span>
                <Info />
              </span>
              <span>{translate("waitingDriver")}</span>
            </div>
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

      <br className="h-16" />

      <div className="flex flex-col gap-2">
        <Button
          disabled={currentEvent !== RiderFlowEvent.WAITING_FOR_DRIVER}
          onClick={handleCancelTrip}
        >
          {translate("cancelTrip")}
        </Button>
      </div>
    </div>
  );
}
