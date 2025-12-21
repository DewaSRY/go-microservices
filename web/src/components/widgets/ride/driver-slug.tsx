import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import { Entity, PackageSlug } from "@/types/common";
import { DriverEvents } from "@/contracts/common";

import { useSocketContext } from "@components/provider/socket-provider";
import useRiderStore from "@/hooks/store/use-ride-store";

import { PackageSlugMeta } from "@components/common/PackagesMeta";

import { cn } from "@/libs/utils";

export default function DriverSettingSlug() {
  const { setPackageSlug } = useUserRideProfile();
  const { sendMessage } = useSocketContext();
  const { currentLocation } = useRiderStore();

  function handlerSelectSlug(packageSLug: PackageSlug) {
    setPackageSlug(packageSLug);

    if (packageSLug !== undefined) {
      sendMessage({
        type: DriverEvents.DRIVER_INIT,
        data: {
          entity: Entity.DRIVER,
          location: currentLocation,
          packageSlug: packageSLug,
        },
      });
    }
  }

  return (
    <div>
      <div className="flex flex-col gap-2">
        {Object.entries(PackageSlugMeta).map(([slug, meta]) => (
          <div
            key={slug}
            className={cn(
              "flex items-center gap-3 sm:gap-4 p-3  transition-all cursor-pointer shadow-lg",
              "sm:p-4 sm:rounded-lg sm:border",
              "hover:border-primary hover:bg-primary/5"
            )}
            onClick={handlerSelectSlug.bind(null, meta.slug)}
          >
            <div className="p-1.5 sm:p-2 bg-gray-100 rounded-lg">
              {meta?.icon}
            </div>
            <div>
              <h3 className="font-medium text-sm sm:text-base">{meta?.name}</h3>
              <p className="text-xs sm:text-sm text-gray-500">
                {meta?.description}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
