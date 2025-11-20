import { Button } from "@components/ui/button";
import { ScrollArea } from "@components/ui/scroll-area";

import useSettingDrawerWidget from "@/hooks/state/useUserSettingDrawer";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";

import UserSettingProfile from "@components/widgets/user-setting-profile";
import UserSettingMode from "@components/widgets/user-setting-mode";
import { cn } from "@/libs/utils";
import { useEffect, useState } from "react";
import UserSettingSlug from "./user-setting-slug";

type userSettingTab = "mode-setting" | "slug-setting" | undefined;

export default function UserSettingSideSheet() {
  const { open, setIsOpen } = useSettingDrawerWidget();
  const { initData, mode, packageSlug } = useUserRideProfile();

  const [currentMode, setTabMode] = useState<userSettingTab>(undefined);

  useEffect(() => {
    initData();
  }, []);

  useEffect(() => {
    console.log("hallo");
    if (mode === undefined) {
      setTabMode("mode-setting");
    } else if (mode !== undefined && packageSlug === undefined) {
      console.log("get hite");
      setTabMode("slug-setting");
    } else {
      setTabMode(undefined);
    }
  }, [mode, packageSlug]);

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
              {currentMode === "mode-setting" && <UserSettingMode />}

              {currentMode === "slug-setting" && <UserSettingSlug />}
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
