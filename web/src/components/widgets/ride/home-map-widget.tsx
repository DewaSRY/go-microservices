import "leaflet/dist/leaflet.css";
import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";
import { MapContainer, Marker, Popup, TileLayer } from "react-leaflet";
import L from "leaflet";
import { useEffect, useMemo, useRef } from "react";
import useRiderStore from "@/hooks/store/use-rider-store";
import MapClickHandler from "../../common/map-click-handler";
import { useSocketContext } from "@components/provider/socket-provider";
import { RoutingControl } from "@components/common/RoutingControl";

import useUserRideProfile from "@/hooks/state/useUserRideProfile";
import { useFlowContext } from "@components/provider/user-flow-provider";
import { RiderFlowEvent } from "@/types/events";
const userMarker = new L.Icon({
  iconUrl:
    "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ed/Map_pin_icon.svg/176px-Map_pin_icon.svg.png",
  iconSize: [40, 40],
  iconAnchor: [20, 40],
});

const driverMarker = new L.Icon({
  iconUrl: "https://www.svgrepo.com/show/25407/car.svg",
  iconSize: [30, 30],
  iconAnchor: [15, 30],
});

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

export default function HomeMapWidget() {
  const mapRef = useRef<L.Map>(null);
  const { destination, setDestination, currentLocation, setLocation } =
    useRiderStore();

  const { currentEvent, isLockDestination, setIsHaveRideRoute } =
    useFlowContext();
  const { routeData, driverActive } = useSocketContext();
  const { mode } = useUserRideProfile();

  const parsedRoute = useMemo(
    () =>
      routeData?.coordinate.map(
        (coord) => [coord.latitude, coord.longitude] as [number, number]
      ),
    [routeData]
  );

  function handlerSelectDestination(e: L.LeafletMouseEvent) {
    if (isLockDestination) return;
    setDestination({
      latitude: e.latlng.lat,
      longitude: e.latlng.lng,
    });
  }

  useEffect(() => {
    setIsHaveRideRoute(!!routeData);
  }, [routeData]);

  useEffect(() => {
    const map = mapRef.current;
    if (map) {
      const timer = setTimeout(() => {
        map.invalidateSize();
      }, 300);

      if ("geolocation" in navigator) {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            const lat = position.coords.latitude;
            const lon = position.coords.longitude;

            setLocation({
              latitude: lat,
              longitude: lon,
            });

            if (mapRef.current) {
              mapRef.current.flyTo([lat, lon], 13, { animate: true });
            }

            setTimeout(() => {
              mapRef.current?.invalidateSize();
            }, 300);
          },
          (error) => {
            console.error("Error getting location:", error);
          }
        );
      }

      return () => {
        clearTimeout(timer);
      };
    }
  }, []);

  return (
    <div className="flex-1 h-full w-full">
      <MapContainer
        center={[currentLocation.latitude, currentLocation.longitude]}
        zoom={13}
        style={{ height: "100%", width: "100%" }}
        ref={mapRef}
      >
        <TileLayer
          url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
          attribution="&copy; <a href='https://www.openstreetmap.org/copyright'>OpenStreetMap</a> contributors &copy; <a href='https://carto.com/'>CARTO</a>"
        />

        {mode === "RIDER" && (
          <Marker
            position={[currentLocation.latitude, currentLocation.longitude]}
            icon={userMarker}
          />
        )}

        {mode === "RIDER" && (
          <div>
            {driverActive.map((x, idx) => (
              <Marker
                key={idx}
                position={[x.coordinate.latitude, x.coordinate.longitude]}
                icon={driverMarker}
              />
            ))}
          </div>
        )}

        {mode === "DRIVER" && (
          <Marker
            position={[currentLocation.latitude, currentLocation.longitude]}
            icon={driverMarker}
          />
        )}

        {destination && mode !== undefined && (
          <Marker
            position={[destination.latitude, destination.longitude]}
            icon={userMarker}
          >
            <Popup></Popup>
          </Marker>
        )}

        {parsedRoute && <RoutingControl route={parsedRoute} />}

        {currentEvent === RiderFlowEvent.TRIP_REQUESTED && (
          <MapClickHandler onClick={handlerSelectDestination} />
        )}
      </MapContainer>
    </div>
  );
}
