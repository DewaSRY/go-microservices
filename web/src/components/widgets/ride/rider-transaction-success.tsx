import TitleBadge from "@/components/ui/title-badge";
import { useTranslations } from "next-intl";
import { Button } from "@components/ui/button";
import { useSocketContext } from "@/components/provider/socket-provider";
import { useFlowContext } from "@/components/provider/user-flow-provider";
import useRiderStore from "@/hooks/store/use-rider-store";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";

export default function RiderTransactionSuccess() {
  const translate = useTranslations("ride.riderTransactionSuccess");
  const { handleReconnect } = useSocketContext();
  const { resetState } = useFlowContext();
  const { reset: riderReset } = useRiderStore();
  const { reset: profileReset } = useUserRideProfile();

  function handleCloseRide() {
    handleReconnect();
    resetState();
    riderReset();
    profileReset();
  }

  return (
    <div>
      <div>
        <TitleBadge>{translate("title")}</TitleBadge>
      </div>
      <br className="h-16" />
      <div className="flex flex-col gap-2">
        <Button onClick={handleCloseRide}>
          {translate("closeTransaction")}
        </Button>
      </div>
    </div>
  );
}
