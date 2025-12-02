import PackageSlugWidget from "@components/widgets/package-slug-widget";

import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import { Entity, PackageSlug } from "@/types/common";
import { DriverEvents } from "@/contracts/common";

import { useSocketContext } from "@components/provider/socket-provider";
import { DriverFlowEvents } from "@/types/events";
import useRiderStore from "@/hooks/store/use-rider-store";

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
      <div>
        <PackageSlugWidget onSelectSlug={handlerSelectSlug} />
      </div>
    </div>
  );
}
