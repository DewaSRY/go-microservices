import { RouteFare, Driver } from "../types/types";
import { TripPreview } from "@/types/dto";
import { DriverList } from "./DriversList";
import { Card } from "./ui/card";
import { Button } from "./ui/button";
import {
  convertMetersToKilometers,
  convertSecondsToMinutes,
} from "../utils/math";
import { Skeleton } from "./ui/skeleton";
import { TripOverviewCard } from "./TripOverviewCard";
import { StripePaymentButton } from "./StripePaymentButton";
import { DriverCard } from "./DriverCard";
import { TripEvents, PaymentEventSessionCreatedData } from "../contracts";
import useRiderStore from "@/hooks/useRiderStore";
import { useTranslations } from "next-intl";

interface TripOverviewProps {
  trip: TripPreview | null;
  onPackageSelect: (carPackage: RouteFare) => void;
  handleAcceptPayment: () => void;
  onCancel: () => void;
}

export const RiderTripOverview = ({
  trip,
  onPackageSelect,
  onCancel,
  handleAcceptPayment,
}: TripOverviewProps) => {
  const riderTripT = useTranslations("riderTrip");
  const commonT = useTranslations("common");

  const assignedDriver = useRiderStore((x) => x.assignedDriver);
  const paymentSession = useRiderStore((x) => x.paymentSession);
  const amountMatchDriver = useRiderStore((x) => x.amountMatchDriver);
  const status = useRiderStore((x) => x.tripStatus);

  if (!trip) {
    return (
      <TripOverviewCard
        title={riderTripT("createTrip.title")}
        description={riderTripT("createTrip.desc")}
      />
    );
  }

  if (status === TripEvents.PaymentSessionCreated && paymentSession) {
    return (
      <TripOverviewCard
        title={riderTripT("tripAccepted.title")}
        description={riderTripT("tripAccepted.desc")}
      >
        <div className="flex flex-col gap-4">
          <DriverCard driver={assignedDriver} />

          <div className="text-sm text-gray-500">
            <p>
              <span>{commonT("amount")}</span>
              <span>:</span>
              <span>
                {paymentSession.amount} {paymentSession.currency}
              </span>
            </p>
            <p>
              <span>{commonT("tripId")}</span>
              <span>:</span>
              <span>{paymentSession.tripID}</span>
            </p>
          </div>
          <StripePaymentButton
            paymentSession={paymentSession}
            handleAcceptPayment={handleAcceptPayment}
          />
        </div>
      </TripOverviewCard>
    );
  }

  if (status === TripEvents.NoDriversFound) {
    return (
      <TripOverviewCard
        title={riderTripT("driverNotFound.title")}
        description={riderTripT("driverNotFound.desc")}
      >
        <Button variant="outline" className="w-full" onClick={onCancel}>
          {commonT("goBack")}
        </Button>
      </TripOverviewCard>
    );
  }

  if (status === TripEvents.DriverAssigned) {
    return (
      <TripOverviewCard
        title={riderTripT("driverAssigned.title")}
        description={riderTripT("driverAssigned.desc")}
      >
        <div className="flex flex-col space-y-3 justify-center items-center mb-4">
          {/* <p>Driver: {trip.id}</p> */}
        </div>
        <Button variant="destructive" className="w-full" onClick={onCancel}>
          {commonT("cancelCurrentTrip")}
        </Button>
      </TripOverviewCard>
    );
  }

  if (status === TripEvents.Completed) {
    return (
      <TripOverviewCard
        title="Trip completed!"
        description="Your trip is completed, thank you for using our service!"
      >
        <Button variant="outline" className="w-full" onClick={onCancel}>
          {commonT("goBack")}
        </Button>
      </TripOverviewCard>
    );
  }

  if (status === TripEvents.Cancelled) {
    return (
      <TripOverviewCard
        title="Trip cancelled!"
        description="Your trip is cancelled, please try again later"
      >
        <Button variant="outline" className="w-full" onClick={onCancel}>
          {commonT("goBack")}
        </Button>
      </TripOverviewCard>
    );
  }

  if (status === TripEvents.Created) {
    return (
      <TripOverviewCard
        title="Looking for a driver"
        description="Your trip is confirmed! We're matching you with a driver, it should not take long."
      >
        <div className="flex flex-col space-y-3 justify-center items-center mb-4">
          <Skeleton className="h-[125px] w-[250px] rounded-xl" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-[250px]" />
            <Skeleton className="h-4 w-[200px]" />
          </div>
        </div>

        <div className="flex flex-col items-center justify-center gap-2">
          {trip?.duration && (
            <h3 className="text-sm font-medium text-gray-700 mb-2">
              Arriving in: {convertSecondsToMinutes(trip?.duration)} at your
              destination ({convertMetersToKilometers(trip?.distance ?? 0)})
            </h3>
          )}

          <Button variant="destructive" className="w-full" onClick={onCancel}>
            {commonT("cancel")}
          </Button>
        </div>
      </TripOverviewCard>
    );
  }

  if (status === TripEvents.TripEventDriversFound) {
    return (
      <TripOverviewCard
        title="Match Driver Found"
        description={`We have ${amountMatchDriver} driver Match Your Needs`}
      >
        <Button variant="outline" className="w-full" onClick={onCancel}>
          {commonT("goBack")}
        </Button>
      </TripOverviewCard>
    );
  }

  if (trip.rideFares && trip.rideFares.length >= 0) {
    return (
      <DriverList
        trip={trip}
        onPackageSelect={onPackageSelect}
        onCancel={onCancel}
      />
    );
  }

  return (
    <Card className="w-full md:max-w-[500px] z-[9999] flex-[0.3]">
      {riderTripT("noTrip.desc")}
    </Card>
  );
};
