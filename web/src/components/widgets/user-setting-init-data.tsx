import { useSocketContext } from "@components/provider/socket-provider";
import { Button } from "../ui/button";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import { RiderEvents } from "@/contracts/common";
import useRiderStore from "@/hooks/store/use-rider-store";

export default function UserSettingInitData() {
  const { sendMessage } = useSocketContext();
  const { mode, packageSlug } = useUserRideProfile();
  const { currentLocation } = useRiderStore();

  function handleInitUserData() {
    if (packageSlug && mode) {
      sendMessage({
        type: RiderEvents.INIT_CONNECTION,
        data: {
          packageSlug: packageSlug,
          entity: mode,
          location: currentLocation,
        },
      });
    }
  }

  return (
    <div>
      <Button onClick={handleInitUserData}>Start Connection</Button>
    </div>
  );
}
