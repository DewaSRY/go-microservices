import { TripOverviewCard } from "../common/TripOverviewCard";
import { Button } from "../ui/button";
import { TripEvents } from "@/_contracts";
import useDriverStore from "@/hooks/useDriverStore";
import { useTranslations } from "next-intl";

export const DriverTripOverview = () => {
  const alertT = useTranslations("alert");
  const driverTripOverViewT = useTranslations("driverTripOverView");

  const driverStore = useDriverStore();
  const trip = useDriverStore((s) => s.trip);
  const tripStatus = useDriverStore((s) => s.tripEvent);

  const handleAcceptTrip = () => {
    const { trip, driver, sendMessage, setTripStatus } = driverStore;
    if (!trip || !trip.id || !driver) {
      alert(alertT("tripIdNotSet"));
      return;
    }

    sendMessage({
      type: TripEvents.DriverTripAccept,
      data: {
        tripID: trip.id,
        riderID: trip.userID,
        driver: driver,
      },
    });

    setTripStatus(TripEvents.DriverTripAccept);
  };

  const handleDeclineTrip = () => {
    const { trip, driver, sendMessage, setTripStatus, resetTripStatus } =
      driverStore;

    if (!trip || !trip.id || !driver) {
      alert(alertT("tripIdNotSet"));
      return;
    }

    sendMessage({
      type: TripEvents.DriverTripDecline,
      data: {
        tripID: trip.id,
        riderID: trip.userID,
        driver: driver,
      },
    });

    setTripStatus(TripEvents.DriverTripDecline);
    resetTripStatus();
  };
  if (!trip) {
    return (
      <TripOverviewCard
        title={driverTripOverViewT("waitingRider.title")}
        description={driverTripOverViewT("waitingRider.desc")}
      />
    );
  }

  if (tripStatus === TripEvents.DriverTripRequest) {
    return (
      <TripOverviewCard
        title={driverTripOverViewT("tripRequest.title")}
        description={driverTripOverViewT("tripRequest.desc")}
      >
        <div className="flex flex-col gap-2">
          <Button onClick={handleAcceptTrip}>
            {driverTripOverViewT("tripRequest.acceptTrip")}
          </Button>
          <Button variant="outline" onClick={handleDeclineTrip}>
            {driverTripOverViewT("tripRequest.declineTrip")}
          </Button>
        </div>
      </TripOverviewCard>
    );
  }

  if (tripStatus === TripEvents.DriverTripAccept) {
    return (
      <TripOverviewCard
        title={driverTripOverViewT("tripAccepted.title")}
        description={driverTripOverViewT("tripAccepted.desc")}
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h3 className="text-lg font-bold">
              {driverTripOverViewT("tripAccepted.tripDetail")}
            </h3>
            <p className="text-sm text-gray-500">
              <span>
                <span>{driverTripOverViewT("tripAccepted.tripId")}</span>
                <span>:</span>
                <span>{trip.id}</span>
              </span>
              <br />
              <span>
                <span>{driverTripOverViewT("tripAccepted.riderId")}</span>
                <span>:</span>
                <span>{trip.userID}</span>
              </span>
            </p>
          </div>
        </div>
      </TripOverviewCard>
    );
  }
  return null;
};
