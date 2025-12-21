import { create } from "zustand";

const initState = {
  open: false,
};

type Action = {
  setIsOpen: (_data: boolean) => void;
};

const useSettingDrawerWidget = create<typeof initState & Action>((set) => {
  function setIsOpen(state: boolean) {
    set((res) => ({
      ...res,
      open: state,
    }));
  }
  return {
    ...initState,
    setIsOpen,
  };
});

export default useSettingDrawerWidget;
