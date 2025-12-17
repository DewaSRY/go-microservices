import useUserMode from "@/hooks/state/useUserRideProfile";
import { Entity } from "@/types/common";
import { useTranslations } from "next-intl";

import { Button } from "@components/ui/button";
import TitleBadge from "@/components/ui/title-badge";

export default function UserSelectMode() {
  const { setMode } = useUserMode();
  const translate = useTranslations("ride.mode");

  return (
    <div className="h-full">
      <TitleBadge>{translate("title")}</TitleBadge>

      <br className="h-16" />

      <div className="flex flex-col gap-2">
        <Button onClick={setMode.bind(null, Entity.DRIVER)}>
          {translate("driverMode")}
        </Button>
        <Button onClick={setMode.bind(null, Entity.RIDER)}>
          {translate("riderMode")}
        </Button>
      </div>

      <br className="h-16" />
    </div>
  );
}
