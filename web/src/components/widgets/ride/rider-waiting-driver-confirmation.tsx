import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";

import TitleBadge from "@/components/ui/title-badge";

export default function RiderWaitingDriverConfirmation() {
  const translate = useTranslations("ride.waitingConfirmation");
  return (
    <div>
      <TitleBadge>{translate("title")}</TitleBadge>
      <br className="h-16" />

      <div className="flex flex-col gap-2">
        <Button>{translate("cancelTrip")}</Button>
      </div>
    </div>
  );
}
