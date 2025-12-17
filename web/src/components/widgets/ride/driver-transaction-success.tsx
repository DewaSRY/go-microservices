import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { useSocketContext } from "@/components/provider/socket-provider";
import { useFlowContext } from "@/components/provider/user-flow-provider";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import TitleBadge from "@/components/ui/title-badge";

export default function DriverTransactionSuccess() {
  const translate = useTranslations("ride.driverTransactionSuccess");

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
        <Button onClick={handleClose}>{translate("closeTransaction")}</Button>
      </div>
      <br className="h-16" />
    </div>
  );
}
