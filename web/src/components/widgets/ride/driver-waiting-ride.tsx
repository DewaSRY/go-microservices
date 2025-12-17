import { useTranslations } from "next-intl";
import TitleBadge from "@/components/ui/title-badge";
import { Button } from "@/components/ui/button";

import { useSocketContext } from "@/components/provider/socket-provider";
import { useFlowContext } from "@/components/provider/user-flow-provider";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";

export default function DriverWaitingRide() {
  const translate = useTranslations("ride.driverWaitingRider");
  const { handleReconnect } = useSocketContext();
  const { resetState } = useFlowContext();
  const { reset: profileReset } = useUserRideProfile();

  function handleClose() {
    handleReconnect();
    resetState();
    profileReset();
  }

  return (
    <div>
      <TitleBadge>{translate("title")}</TitleBadge>
      <br className="h-16" />
      <div className="flex flex-col gap-2">
        <Button onClick={handleClose}>{translate("close")}</Button>
      </div>
      <br className="h-16" />
    </div>
  );
}
