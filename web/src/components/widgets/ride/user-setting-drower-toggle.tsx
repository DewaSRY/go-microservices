import useUserSettingDrawer from "@hooks/state/useUserSettingDrawer";
import { Button } from "@components/ui/button";
import { Fragment } from "react";
import { Menu } from "lucide-react";

export default function userSettingMenuButtonWidget() {
  const { open, setIsOpen } = useUserSettingDrawer();
  return (
    <Fragment>
      {!open && (
        <div>
          <Button onClick={setIsOpen.bind(null, true)} size="icon">
            <Menu />
          </Button>
        </div>
      )}
    </Fragment>
  );
}
