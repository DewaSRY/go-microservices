"use client";
import HomeMapWidget from "../widgets/ride/home-map-widget";
import { HomeNavigationWidget } from "@/components/widgets/home-navigation";
import SocketProvider from "@components/provider/socket-provider";
import UserSettingDrawer from "@/components/widgets/ride/user-setting-drawer";
import UserSettingMenuButtonWidget from "@/components/widgets/ride/user-setting-drower-toggle";
import UserFlowProvider from "@components/provider/user-flow-provider";

export default function HeroPage() {
  return (
    <main className="relative">
      <SocketProvider>
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
      </SocketProvider>
    </main>
  );
}
