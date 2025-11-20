import useUserMode from "@/hooks/state/useUserRideProfile";

import { Button } from "@components/ui/button";

export default function userSettingMode() {
  const { setMode } = useUserMode();

  return (
    <div className="h-full">
      <div className="mb-8"></div>
      <div className="flex flex-col gap-2">
        <Button onClick={setMode.bind(null, "DRIVER")}>Driver</Button>
        <Button onClick={setMode.bind(null, "RIDER")}>Rider</Button>
      </div>
    </div>
  );
}
