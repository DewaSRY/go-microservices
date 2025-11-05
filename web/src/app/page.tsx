"use client";

// Assets
import "leaflet/dist/leaflet.css";
// Fix for default marker icon
import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";
import dynamic from "next/dynamic";
import { Button } from "../components/ui/button";
import { useState, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { CarPackageSlug } from "../types/types";
import { DriverPackageSelector } from "@/components/driver/DriverPackageSelector";
import { useTranslations } from "next-intl";

// Dynamically import components that use Leaflet
const DriverMapWidget = dynamic(
  () =>
    import("@/components/driver/DriverMapWidget").then((mod) => mod.DriverMap),
  { ssr: false }
);

const RiderMapWidget = dynamic(
  () => import("../components/rider/RiderMapWidget"),
  {
    ssr: false,
  }
);

// Initialize Leaflet icon only on client side
if (typeof window !== "undefined") {
  import("leaflet").then((L) => {
    const DefaultIcon = L.default.icon({
      iconUrl: icon.src,
      shadowUrl: iconShadow.src,
      iconSize: [25, 41],
      iconAnchor: [12, 41],
    });
    L.default.Marker.prototype.options.icon = DefaultIcon;
  });
}

function HomeContent() {
  const t = useTranslations("HomePage");

  const router = useRouter();
  const searchParams = useSearchParams();
  const payment = searchParams.get("payment");

  const [userType, setUserType] = useState<"driver" | "rider" | null>(null);
  const [packageSlug, setPackageSlug] = useState<CarPackageSlug | null>(null);

  const handleClick = (userType: "driver" | "rider") => {
    setUserType(userType);
  };

  const handleBackToHome = () => {
    router.push("/");
  };

  if (payment === "success") {
    return (
      <main className="min-h-screen bg-gradient-to-b from-white to-gray-50">
        <div className="flex flex-col items-center justify-center h-screen gap-6 px-4">
          <div className="bg-white p-8 rounded-2xl shadow-lg text-center max-w-md w-full">
            <div className="mb-6">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg
                  className="w-8 h-8 text-green-500"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900">
                {t("successWidget.title")}
              </h1>
              <p className="text-gray-600 mt-2">{t("successWidget.desc")}</p>
            </div>
            <Button
              className="w-full text-lg py-6"
              variant="outline"
              onClick={handleBackToHome}
            >
              {t("successWidget.returnHome")}
            </Button>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-gradient-to-b from-white to-gray-50">
      {userType === null && (
        <div className="flex flex-col items-center justify-center h-screen gap-6 px-4">
          <div className="bg-white p-8 rounded-2xl shadow-lg text-center max-w-md w-full">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">
              {t("hero.title")}
            </h2>
            <p className="text-gray-600 mb-8">{t("hero.desc")}</p>
            <div className="space-y-4">
              <Button
                className="w-full text-lg py-6 bg-primary hover:bg-primary/90"
                onClick={handleClick.bind(null, "rider")}
              >
                {t("hero.wantToRide")}
              </Button>
              <Button
                className="w-full text-lg py-6"
                variant="outline"
                onClick={handleClick.bind(null, "driver")}
              >
                {t("hero.wantToDrive")}
              </Button>
            </div>
          </div>
        </div>
      )}

      {userType === "driver" && packageSlug && (
        <DriverMapWidget packageSlug={packageSlug} />
      )}

      {userType === "driver" && !packageSlug && (
        <DriverPackageSelector onSelect={setPackageSlug} />
      )}

      {userType === "rider" && <RiderMapWidget />}
    </main>
  );
}

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
      <HomeContent />
    </Suspense>
  );
}
