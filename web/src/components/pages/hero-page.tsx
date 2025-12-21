"use client";
import HomeMapWidget from "../widgets/ride/home-map-widget";
import { HomeNavigationWidget } from "@/components/widgets/home-navigation";
import SocketProvider from "@components/provider/socket-provider";
import UserLocationProvider from "@components/provider/user-location-provider";
import UserSettingDrawer from "@/components/widgets/ride/user-setting-drawer";
import UserSettingMenuButtonWidget from "@/components/widgets/ride/user-setting-drawer-toggle";
import UserFlowProvider from "@components/provider/user-flow-provider";
import { useEffect } from "react";
import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import useRideStore from "@/hooks/store/use-ride-store";

export default function HeroPage() {
  const { initData: riderInitData } = useUserRideProfile();
  const { initData: rideInitData } = useRideStore();

  useEffect(() => {
    rideInitData();
    riderInitData();
  }, []);

  return (
    <main className="relative">
      <SocketProvider>
        <UserLocationProvider>
          <UserFlowProvider>
            <div className="absolute top-2 left-1/2 -translate-x-1/2  z-[1000]">
              <HomeNavigationWidget />
            </div>
            <div className="relative flex flex-col md:flex-row h-screen z-0">
              <HomeMapWidget />
            </div>

            <div className="z-[1000] absolute top-[12%] right-2">
              <UserSettingMenuButtonWidget />
            </div>

            <div className="z-[1000]">
              <UserSettingDrawer />
            </div>
          </UserFlowProvider>
        </UserLocationProvider>
      </SocketProvider>
    </main>
  );
}
