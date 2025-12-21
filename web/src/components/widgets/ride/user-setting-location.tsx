import { useTranslations } from "next-intl";

import TitleBadge from "@/components/ui/title-badge";
import { Button } from "@/components/ui/button";
import { useUserLocationContext } from "@/components/provider/user-location-provider";

export default function UserSettingSLocation() {
  const translate = useTranslations("ride.settingLocation");
  const { requestUserLocation } = useUserLocationContext();

  return (
    <div>
      <TitleBadge>{translate("title")}</TitleBadge>
      <br className="h-16" />

      <div className="flex flex-col gap-2">
        <Button onClick={requestUserLocation}>
          {translate("enableLocation")}
        </Button>
      </div>
      <br className="h-16" />
    </div>
  );
}
