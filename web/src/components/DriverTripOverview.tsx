import { TripOverviewCard } from "./TripOverviewCard";
import { Button } from "./ui/button";
import { TripEvents } from "../contracts";
import useDriverStore from "@/hooks/useDriverStore";

export const DriverTripOverview = () => {
  const driverStore = useDriverStore();
  const trip = useDriverStore((s) => s.trip);
  const tripStatus = useDriverStore((s) => s.tripEvent);

  const handleAcceptTrip = () => {
    const { trip, driver, sendMessage, setTripStatus } = driverStore;
    if (!trip || !trip.id || !driver) {
      alert("No trip ID found or driver is not set");
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
      alert("No trip ID found or driver is not set");
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
        title="Waiting for a rider..."
        description="Waiting for a rider to request a trip..."
      />
    );
  }

  if (tripStatus === TripEvents.DriverTripRequest) {
    return (
      <TripOverviewCard
        title="Trip request received!"
        description="A trip has been requested, check the route and accept the trip if you can take it."
      >
        <div className="flex flex-col gap-2">
          <Button onClick={handleAcceptTrip}>Accept trip</Button>
          <Button variant="outline" onClick={handleDeclineTrip}>
            Decline trip
          </Button>
        </div>
      </TripOverviewCard>
    );
  }

  if (tripStatus === TripEvents.DriverTripAccept) {
    return (
      <TripOverviewCard
        title="All set!"
        description="You can now start the trip"
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h3 className="text-lg font-bold">Trip details</h3>
            <p className="text-sm text-gray-500">
              Trip ID: {trip.id}
              <br />
              Rider ID: {trip.userID}
            </p>
          </div>
        </div>
      </TripOverviewCard>
    );
  }
  return null;
};
