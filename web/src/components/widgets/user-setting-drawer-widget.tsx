import { Button } from "@components/ui/button";
import { useState } from "react";
import { ScrollArea, ScrollBar } from "@components/ui/scroll-area";

import useSettingDrawerWidget from "@/hooks/state/useUserSettingDrawer";
import { cn } from "@/libs/utils";

export default function UserSettingSideSheet() {
  const { open, setIsOpen } = useSettingDrawerWidget();

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
          <header className="p-4 border-b">
            <h2 className="text-lg font-semibold">User Settings</h2>
            <p className="text-sm text-muted-foreground">
              Manage your account preferences.
            </p>
          </header>

          <div className="p-4 flex-1 overflow-auto">
            {/* Your content here */}
            Settings content...
          </div>

          <footer className="p-4 border-t flex justify-between">
            <Button>Save</Button>
            <Button variant="outline" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
          </footer>
        </ScrollArea>
      </div>
    </>
  );
}
