import { create } from "zustand";
import { PackageSlug, Entity } from "@/types/common";

type modeType = Entity | undefined;

const initState = {
  mode: undefined as modeType,
  packageSlug: undefined as PackageSlug | undefined,
};

type Action = {
  setMode: (_data: modeType) => void;
  setPackageSlug: (_data: PackageSlug | undefined) => void;
  initData: () => void;
};

const useUserRideProfile = create<typeof initState & Action>((set) => {
  function setMode(data: Entity | undefined) {
    set((res) => {
      return {
        ...res,
        mode: data,
      };
    });

    if (data !== undefined) {
      // localStorage.setItem("mode", data);
    } else {
      localStorage.removeItem("mode");
    }
  }

  function setPackageSlug(data: PackageSlug | undefined) {
    set((res) => {
      return {
        ...res,
        packageSlug: data,
      };
    });

    if (data !== undefined) {
      // localStorage.setItem("package-slug", data);
    } else {
      localStorage.removeItem("package-slug");
    }
  }

  function initData() {
    const storedMode = localStorage.getItem("mode");
    const storedPackageSLug = localStorage.getItem("package-slug");

    let mode: modeType = undefined;
    let packageSlug: PackageSlug;

    if (
      storedMode &&
      Object.values(Entity)
        .map((x) => x.toString())
        .includes(storedMode)
    ) {
      mode = storedMode as Entity;
    }

    if (
      storedPackageSLug &&
      Object.values(PackageSlug)
        .map((x) => x.toString())
        .includes(storedPackageSLug)
    ) {
      packageSlug = storedPackageSLug as PackageSlug;
    }

    set((res) => {
      return {
        ...res,
        mode: mode,
        packageSlug: packageSlug,
      };
    });
  }

  return {
    ...initState,
    setMode,
    initData,
    setPackageSlug,
  };
});

export default useUserRideProfile;
