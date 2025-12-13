import { DriverEvents } from "@/contracts/common";
import { useSocketContext } from "../provider/socket-provider";
import { Button } from "../ui/button";

export default function DriverTransaction() {
  const { sendMessage, transactionId } = useSocketContext();

  function handleAcceptedTransaction() {
    if (transactionId !== undefined) {
      sendMessage({
        type: DriverEvents.DRIVER_ACCEPTED_TRANSACTION,
        data: {
          transactionId: transactionId,
        },
      });
    }
  }

  return (
    <div>
      <div>Driver Find rider</div>

      <div>
        <Button onClick={handleAcceptedTransaction}>accepted</Button>
      </div>
    </div>
  );
}
