import { HTMLAttributes } from "react";
import { Button } from "../../ui/button";
import { useSocketContext } from "@/components/provider/socket-provider";
import { RiderEvents } from "@/contracts/common";
import { PackageSlugMeta } from "@/components/common/PackagesMeta";
import { PackageSlug } from "@/types/common";

export interface DriverActiveCardProps extends HTMLAttributes<HTMLDivElement> {
  packageSlug: PackageSlug;
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
    <div className="flex items-center border rounded-lg shadow-md px-2 py-2">
      <div className="flex-1">
        <h3 className="font-medium text-gray-900 text-lg uppercase">{title}</h3>
        <div className="flex gap-2">
          <div>{PackageSlugMeta[packageSlug].icon}</div>
          <span className="font-extralight text-gray-500 text-md ">
            {packageSlug}
          </span>
        </div>
      </div>
      <div>
        <Button onClick={handleSelectActiveDriver} size="sm">
          Select
        </Button>
      </div>
    </div>
  );
}
