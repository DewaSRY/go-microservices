import { HTMLAttributes } from "react";
import { Button } from "../ui/button";
import { useSocketContext } from "../provider/socket-provider";
import { RiderEvents } from "@/contracts/common";

export interface DriverActiveCardProps extends HTMLAttributes<HTMLDivElement> {
  packageSlug: string;
  title: string;
  driverId: string;
}

export default function DriverActiveCard({
  packageSlug,
  title,
  driverId,
}: DriverActiveCardProps) {
  const { sendMessage } = useSocketContext();

  function handleSelectActiveDriver() {
    sendMessage({
      type: RiderEvents.RIDER_CREATE_TRANSACTION,
      data: {
        driverId: driverId,
      },
    });
  }

  return (
    <div className="mb-2">
      <div>
        <h3 className="font-medium text-gray-600">{title}</h3>
        <span className="font-extralight text-gray-500">{packageSlug}</span>
      </div>
      <div>
        <Button onClick={handleSelectActiveDriver}>Select</Button>
      </div>
    </div>
  );
}
