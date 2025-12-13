"use client";

import { Suspense } from "react";
import dynamic from "next/dynamic";

const HeroPage = dynamic(
  () => import("@/components/pages/hero-page").then((mod) => mod.default),
  { ssr: false }
);

export default function Home() {
  return (
    <Suspense
      fallback={
        <main className="min-h-screen bg-gradient-to-b from-white to-gray-50">
          <div className="flex flex-col items-center justify-center h-screen gap-4">
            <div className="bg-white p-8 rounded-2xl shadow-lg text-center max-w-md w-full">
              <div className="animate-pulse flex flex-col items-center">
                <div className="h-8 w-32 bg-gray-200 rounded mb-4"></div>
                <div className="h-4 w-48 bg-gray-100 rounded"></div>
              </div>
            </div>
          </div>
        </main>
      }
    >
      <HeroPage />
    </Suspense>
  );
}
