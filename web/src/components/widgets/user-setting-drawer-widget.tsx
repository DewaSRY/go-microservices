import { Button } from "@components/ui/button";
import { ScrollArea } from "@components/ui/scroll-area";

import useSettingDrawerWidget from "@/hooks/state/useUserSettingDrawer";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import UserSettingProfile from "@components/widgets/user-setting-profile";
import UserSelectMode from "@/components/widgets/user-select-mode";
import { cn } from "@/libs/utils";
import { useEffect } from "react";

import { useFlowContext } from "@components/provider/user-flow-provider";
import { RiderFlowEvent } from "@/types/events";
import UserSettingSelectedTrip from "./user-setting-selected-trip";

export default function UserSettingSideSheet() {
  const { currentEvent } = useFlowContext();

  const { open, setIsOpen } = useSettingDrawerWidget();
  const { initData, mode } = useUserRideProfile();

  useEffect(() => {
    initData();
  }, []);

  return (
    <>
      {/* Your trigger button anywhere */}
      <Button onClick={() => setIsOpen(true)}>Open Settings</Button>

      {/* Slide-in panel */}
      <div
        className={cn(
          "fixed top-[10%] right-4 h-[80vh] w-[360px] bg-white border-l rounded-md shadow-xl",
          "translate-x-[120%] transition-transform duration-300 ease-out",
          open && " translate-x-0"
        )}
      >
        <ScrollArea>
          <div className="">
            <header className="p-4 border-b">
              <UserSettingProfile />
            </header>

            <div className="p-4 flex-1 overflow-auto grow h-[80%] ">
              {mode === undefined && <UserSelectMode />}

              {/* {currentEvent === RiderFlowEvent.RIDER_INIT_CONNECTION && (
                <UserSettingSlug />
              )} */}

              {(currentEvent === RiderFlowEvent.TRIP_REQUESTED ||
                currentEvent === RiderFlowEvent.WAITING_FOR_DRIVER) && (
                <UserSettingSelectedTrip />
              )}

              {/* {currentMode === "slug-setting" && <UserSettingSlug />}

              {currentMode === "start-connection-setting" && (
                <UserSettingInitData />
              )}

              {currentMode === "init-data-success" && (
                <RiderSelectDestination />
              )}

              {currentMode === undefined && <p>this is undefine</p>} */}
            </div>

            <footer className="p-4 border-t flex justify-between">
              <hr />
            </footer>
          </div>
        </ScrollArea>
      </div>
    </>
  );
}
