import { useSocketContext } from "@components/provider/socket-provider";
import { Button } from "../ui/button";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import { Entity } from "@/types/common";
import { RiderEvents } from "@/contracts/common";
export default function UserSettingInitData() {
  const { sendMessage } = useSocketContext();
  const { mode, packageSlug } = useUserRideProfile();

  function handleInitUserData() {
    if (packageSlug && mode) {
      sendMessage({
        type: RiderEvents.INIT_CONNECTION,
        data: {
          packageSlug: packageSlug,
          entity: mode,
          location: {
            latitude: 0,
            longitude: 0,
          },
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
