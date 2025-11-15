"use client";
import HomeMapWidget from "../widgets/HomeMapWidget";
import { HomeNavigationWidget } from "../widgets/HomeNavigationWidget";

export default function HeroPage() {
  return (
    <main className="relative">
      <div className="absolute top-2 left-1/2 -translate-x-1/2  z-[1000]">
        <HomeNavigationWidget />
      </div>
      <div className="relative flex flex-col md:flex-row h-screen">
        <HomeMapWidget />
      </div>
    </main>
  );
}
