import {
  PropsWithChildren,
  useContext,
  HTMLAttributes,
  createContext,
  useEffect,
  use,
  useState,
} from "react";
import { RiderFlowEvent, DriverFlowEvents } from "@/types/events";
import useRiderStore from "@/hooks/store/use-ride-store";
import { LocationPermissionState } from "@/types/common";

type UserLocationContextType = {
  requestUserLocation: () => Promise<void>;
  locationPermission: LocationPermissionState;
};

const UserLocationProvider = createContext<UserLocationContextType | undefined>(
  undefined
);

UserLocationProvider.displayName = "user-location";

interface ProviderPops
  extends HTMLAttributes<HTMLDivElement>,
    PropsWithChildren {}

export default function Provider({ children }: ProviderPops) {
  const { setLocation } = useRiderStore();

  const [locationPermission, setLocationPermission] =
    useState<LocationPermissionState>("unsupported");

  useEffect(() => {
    if ("geolocation" in navigator) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const lat = position.coords.latitude;
          const lon = position.coords.longitude;

          setLocation({
            latitude: lat,
            longitude: lon,
          });
        },
        (error) => {
          console.error("Error getting location:", error);
        }
      );
    }
  }, []);

  useEffect(() => {
    let permissionStatus: PermissionStatus | null = null;
    async function initPermission() {
      if (!("permissions" in navigator)) {
        setLocation;
      }

      try {
        permissionStatus = await navigator.permissions.query({
          name: "geolocation",
        });
        setLocationPermission(
          permissionStatus.state as LocationPermissionState
        );

        permissionStatus.onchange = () => {
          setLocationPermission(
            permissionStatus?.state as LocationPermissionState
          );
        };
      } catch (err) {
        console.error("Error getting location permission:", err);
        setLocationPermission("unsupported");
      }
    }

    initPermission();

    return () => {
      if (permissionStatus) {
        permissionStatus.onchange = null;
      }
    };
  }, []);

  function requestUserLocation(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      if (!("geolocation" in navigator)) {
        reject(new Error("Geolocation not supported"));
      } else {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            const lat = position.coords.latitude;
            const lon = position.coords.longitude;

            setLocation({
              latitude: lat,
              longitude: lon,
            });
            resolve();
          },
          (error) => {
            console.error("Error getting location:", error);
            reject(error);
          },
          {
            enableHighAccuracy: false,
            timeout: 10000,
          }
        );
      }
    });
  }

  return (
    <UserLocationProvider.Provider
      value={{
        requestUserLocation: requestUserLocation,
        locationPermission: locationPermission,
      }}
    >
      {children}
    </UserLocationProvider.Provider>
  );
}

export function useUserLocationContext() {
  const context = useContext(UserLocationProvider);

  if (!context) throw Error("hook_use_outside_user_flow_the_context");

  return context;
}
