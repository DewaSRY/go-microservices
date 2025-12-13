import { Button } from "@components/ui/button";
import { ScrollArea } from "@components/ui/scroll-area";

import useSettingDrawerWidget from "@/hooks/state/useUserSettingDrawer";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import UserSettingProfile from "@components/widgets/user-setting-profile";
import UserSelectMode from "@/components/widgets/user-select-mode";
import { cn } from "@/libs/utils";
import { useEffect } from "react";

import { useFlowContext } from "@components/provider/user-flow-provider";
import { DriverFlowEvents, RiderFlowEvent } from "@/types/events";
import RiderListingDriver from "./driver-setting-slug";
import RiderSettingWaitingDriver from "./rider-listing-driver";
import RiderCreateTrip from "./rider-create-trip copy";
import DriverTransaction from "./driver-transaction";
import DriverWaitingRide from "./driver-waiting-ride";
import RiderWaitingDriverConfirmation from "./rider-waiting-driver-confirmation";
import RiderTransactionSuccess from "./rider-transaction-success";
import DriverTransactionSuccess from "./driver-transaction-success";

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
      <Button onClick={setIsOpen.bind(null, true)}>Open Settings</Button>

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

              {currentEvent === RiderFlowEvent.TRIP_REQUESTED && (
                <RiderCreateTrip />
              )}

              {currentEvent === RiderFlowEvent.WAITING_FOR_DRIVER && (
                <RiderSettingWaitingDriver />
              )}

              {currentEvent === DriverFlowEvents.DRIVER_INIT_CONN && (
                <RiderListingDriver />
              )}

              {currentEvent ===
                RiderFlowEvent.RIDER_WAITING_DRIVER_CONFIRMATION && (
                <RiderWaitingDriverConfirmation />
              )}

              {currentEvent === RiderFlowEvent.RIDER_TRANSACTION_SUCCESS && (
                <RiderTransactionSuccess />
              )}

              {currentEvent === DriverFlowEvents.DRIVER_WAITING_FOR_RIDER && (
                <DriverWaitingRide />
              )}

              {currentEvent === DriverFlowEvents.RIDER_CREATE_TRANSACTION && (
                <DriverTransaction />
              )}

              {currentEvent === DriverFlowEvents.DRIVER_TRANSACTION_SUCCESS && (
                <DriverTransactionSuccess />
              )}
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
