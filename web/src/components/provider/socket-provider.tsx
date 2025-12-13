import { UserWsRequest } from "@/contracts/ws-request";
import { RiderWsResponse } from "@/contracts/ws-response";
import { useCallback, useRef, useState } from "react";
import { useEffect } from "react";
import {
  PropsWithChildren,
  useContext,
  HTMLAttributes,
  createContext,
} from "react";

import { RideEvents, RiderEvents } from "@/contracts/common";

import { WEBSOCKET_URL } from "@constants/environment";
import { DriverActiveRecord, RouteData } from "@/types/common";

type ConnectionState = RiderEvents | undefined;

const SocketProviderContext = createContext({
  sendMessage: (_data: UserWsRequest) => {},
  isConnected: false,
  connectionState: undefined as ConnectionState,
  isLoading: false,
  routeData: undefined as RouteData | undefined,
  resetRoute: () => {},
  driverActive: [] as DriverActiveRecord[],
  transactionId: undefined as string | undefined,
  isTransactionAccepted: undefined as boolean | undefined,
});

SocketProviderContext.displayName = "socket-provider";

interface ProviderPops
  extends HTMLAttributes<HTMLDivElement>,
    PropsWithChildren {
  reconnect?: boolean;
  reconnectInterval?: number;
}

export default function Provider({
  children,
  reconnect = true,
  reconnectInterval = 3000,
}: ProviderPops) {
  const [isLoading, setIsLoading] = useState(false);
  const [currentState, setCurrentState] = useState<ConnectionState>(undefined);
  const [isConnected, setIsConnected] = useState(false);
  const [route, setRoute] = useState<RouteData | undefined>(undefined);
  const [driverActive, setDriverActive] = useState<DriverActiveRecord[]>([]);
  const [transactionId, setTransactionId] = useState<string | undefined>(
    undefined
  );
  const [isTransactionAccepted, setIsTransactionAccepted] = useState<
    boolean | undefined
  >(undefined);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const socketRef = useRef<WebSocket | null>(null);

  const connect = useCallback(() => {
    try {
      const ws = new WebSocket(`${WEBSOCKET_URL}/connect`);
      ws.onmessage = (event) => {
        const message = JSON.parse(event.data) as RiderWsResponse;

        console.log(message);
        switch (message.type) {
          case RiderEvents.CONNECTION_SUCCESS:
            setCurrentState(message.type);
            break;
          case RideEvents.ROUTE_FOUND_RESPONSE:
            setRoute(message.data);
            break;
          case RideEvents.DRIVER_ACTIVE_RESPONSE:
            console.log(message);
            setDriverActive(() => {
              return message.data;
            });
            break;
          case RideEvents.RIDER_CREATE_TRANSACTION_RESPONSE:
            setTransactionId(message.data.transactionId);
            break;
          case RideEvents.TRANSACTION_ACCEPTED_RESPONSE:
            setIsTransactionAccepted(true);
            break;
          default:
            setCurrentState(undefined);
        }

        setIsLoading((prev) => {
          if (prev == true) {
            return false;
          }
          return prev;
        });
      };

      ws.onclose = () => {
        setIsConnected(false);
        if (reconnect) {
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        }
      };

      socketRef.current = ws;
    } catch (e) {
      console.log("failed_to_create_websocket:", e);
    }
  }, []);

  useEffect(() => {
    setRoute(undefined);
    connect();
  }, []);

  const sendMessage = useCallback((data: UserWsRequest) => {
    const ws = socketRef.current;
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data));
    } else {
      console.warn("websocket_is_not_connect");
    }
  }, []);

  function resetRoute() {
    setRoute(undefined);
  }

  return (
    <SocketProviderContext.Provider
      value={{
        sendMessage: sendMessage,
        isConnected,
        connectionState: currentState,
        isLoading,
        routeData: route,
        resetRoute,
        driverActive: driverActive,
        transactionId,
        isTransactionAccepted,
      }}
    >
      {children}
    </SocketProviderContext.Provider>
  );
}

export function useSocketContext() {
  const context = useContext(SocketProviderContext);

  if (!context) throw Error("hook_use_outside_the_context");

  return context;
}
