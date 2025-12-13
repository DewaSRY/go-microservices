import useUserMode from "@/hooks/state/useUserRideProfile";
import { Entity } from "@/types/common";

import { Button } from "@components/ui/button";

export default function UserSelectMode() {
  const { setMode } = useUserMode();

  return (
    <div className="h-full">
      <div className="mb-8"></div>
      <div className="flex flex-col gap-2">
        <Button onClick={setMode.bind(null, Entity.DRIVER)}>Driver</Button>
        <Button onClick={setMode.bind(null, Entity.RIDER)}>Rider</Button>
      </div>
    </div>
  );
}
