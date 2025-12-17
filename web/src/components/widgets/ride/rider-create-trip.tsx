import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import useRiderStore from "@/hooks/store/use-rider-store";
import { useSocketContext } from "@/components/provider/socket-provider";
import { RiderEvents } from "@/contracts/common";
import TitleBadge from "@/components/ui/title-badge";

export default function RiderCreateTrip() {
  const { currentLocation, destination } = useRiderStore();
  const { sendMessage } = useSocketContext();
  const translate = useTranslations("ride.rideCreateTrip");

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
      <TitleBadge>{translate("title")}</TitleBadge>
      <br className="h-16" />

      <div className="flex flex-col gap-2">
        <Button disabled={destination === undefined} onClick={handleStartTrip}>
          {translate("startTrip")}
        </Button>
      </div>
      <br className="h-16" />
    </div>
  );
}
