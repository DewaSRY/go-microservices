import { RiderWsMessage } from "@/contracts/rider-connection";
import { useCallback, useRef, useState } from "react";
import { useEffect } from "react";
import {
  PropsWithChildren,
  useContext,
  HTMLAttributes,
  createContext,
} from "react";

import { WEBSOCKET_URL } from "@constants/environment";
const SocketProviderContext = createContext({
  sendMessage: (_data: RiderWsMessage) => {},
  isConnected: false,
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
  const [isConnected, setIsConnected] = useState(false);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const socketRef = useRef<WebSocket | null>(null);

  const connect = useCallback(() => {
    try {
      const ws = new WebSocket(`${WEBSOCKET_URL}/connection`);

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
  }, [reconnect, reconnectInterval]);

  useEffect(() => {
    const socket = new WebSocket(`${WEBSOCKET_URL}/connect`);

    socketRef.current = socket;
  });

  const sendMessage = useCallback((data: RiderWsMessage) => {
    const ws = socketRef.current;
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data));
    } else {
      console.warn("websocket_is_not_connect");
    }
  }, []);

  return (
    <SocketProviderContext.Provider
      value={{
        sendMessage: sendMessage,
        isConnected,
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
