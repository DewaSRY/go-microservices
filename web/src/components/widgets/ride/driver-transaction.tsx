import { DriverEvents } from "@/contracts/common";
import { useSocketContext } from "@/components/provider/socket-provider";
import { Button } from "@/components/ui/button";
import { useTranslations } from "next-intl";
import TitleBadge from "@/components/ui/title-badge";

export default function DriverTransaction() {
  const translate = useTranslations("ride.driverGetRiderTransaction");

  const { sendMessage, transactionId } = useSocketContext();

  function handleAcceptedTransaction() {
    if (transactionId !== undefined) {
      sendMessage({
        type: DriverEvents.DRIVER_ACCEPTED_TRANSACTION,
        data: {
          transactionId: transactionId,
        },
      });
    }
  }

  return (
    <div>
      <TitleBadge>{translate("title")}</TitleBadge>
      <br className="h-16" />

      <div className="flex flex-col gap-2">
        <Button onClick={handleAcceptedTransaction}>
          {translate("acceptTransaction")}
        </Button>
      </div>

      <br className="h-16" />
    </div>
  );
}
