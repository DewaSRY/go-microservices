import useUserSettingDrawer from "@hooks/state/useUserSettingDrawer";
import { Button } from "@components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@components/ui/avatar";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";

export default function userSettingProfile() {
  const { setIsOpen } = useUserSettingDrawer();
  const { mode, setMode, setPackageSlug, packageSlug } = useUserRideProfile();

  function handleChangeData() {
    setMode(undefined);
    setPackageSlug(undefined);
  }

  return (
    <div className="flex justify-between items-center">
      <div>
        <Button onClick={setIsOpen.bind(null, false)}>Close</Button>
      </div>

      <div className="flex justify-end gap-2">
        <div className="flex flex-col">
          <span className="text-lg text-right text-gray-900">unknow</span>

          <div className="text-xs text-right text-gray-700 flex gap-1 items-center">
            <div className="flex gap-1">
              {mode === undefined ? (
                <span>(undefined)</span>
              ) : (
                <div className="flex gap-2 items-center">
                  <span>{mode}</span>
                </div>
              )}

              {packageSlug === undefined ? (
                <span>(undefined)</span>
              ) : (
                <div className="flex gap-2 items-center">
                  <span>{packageSlug}</span>
                </div>
              )}

              {mode && <button onClick={handleChangeData}>change</button>}
            </div>
            <span></span>
          </div>
        </div>

        <div>
          <Avatar>
            <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </div>
      </div>
    </div>
  );
}
