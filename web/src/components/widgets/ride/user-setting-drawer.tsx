import { ScrollArea } from "@components/ui/scroll-area";

import useSettingDrawerWidget from "@/hooks/state/useUserSettingDrawer";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import UserSettingProfile from "@/components/widgets/ride/user-profile";
import UserSelectMode from "@/components/widgets/ride/user-select-mode";
import { cn } from "@/libs/utils";
import { useEffect } from "react";

import { useFlowContext } from "@components/provider/user-flow-provider";
import { DriverFlowEvents, RiderFlowEvent } from "@/types/events";
import RiderListingDriver from "./driver-slug";
import RiderSettingWaitingDriver from "./rider-listing-driver";
import RiderCreateTrip from "./rider-create-trip";
import DriverTransaction from "./driver-transaction";
import DriverWaitingRide from "./driver-waiting-ride";
import RiderWaitingDriverConfirmation from "./rider-waiting-driver-confirmation";
import RiderTransactionSuccess from "./rider-transaction-success";
import DriverTransactionSuccess from "./driver-transaction-success";

export default function UserSettingSideSheet() {
  const { currentEvent } = useFlowContext();
  const { open } = useSettingDrawerWidget();
  const { initData, mode } = useUserRideProfile();

  useEffect(() => {
    initData();
  }, []);

  return (
    <div
      className={cn(
        "fixed top-[10%] right-4 h-[80vh] w-[360px] bg-white border-l rounded-md shadow-xl",
        "translate-x-[120%] transition-transform duration-300 ease-out",
        open && " translate-x-0"
      )}
    >
      <ScrollArea>
        <div className="flex flex-col h-full">
          <header className="p-4 border-b">
            <UserSettingProfile />
          </header>

          <div className="p-4 flex-1 overflow-auto grow h-5/6">
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
        </div>
      </ScrollArea>
    </div>
  );
}
