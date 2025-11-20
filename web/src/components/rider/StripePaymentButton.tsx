import { useTranslations } from "next-intl";
import { PaymentEventSessionCreatedData } from "@/_contracts";
import { Button } from "@/components/ui/button";
import { loadStripe } from "@stripe/stripe-js";

interface StripePaymentButtonProps {
  paymentSession: PaymentEventSessionCreatedData;
  isLoading?: boolean;
  handleAcceptPayment: () => void;
}

// Initialize Stripe
const stripePromise = loadStripe(
  process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!
);

export const StripePaymentButton = ({
  paymentSession,
  handleAcceptPayment,
  isLoading = false,
}: StripePaymentButtonProps) => {
  const paymentT = useTranslations("payment");

  const handlePayment = async () => {
    const stripe = await stripePromise;

    if (!stripe) {
      console.error("Stripe failed to load");
      return;
    }

    handleAcceptPayment();
  };

  if (!process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY) {
    return (
      <Button disabled className="w-full bg-red-500 text-white">
        {paymentT("warning.noStripApiKey")}
      </Button>
    );
  }

  return (
    <Button onClick={handlePayment} disabled={isLoading} className="w-full">
      {isLoading
        ? "Loading..."
        : `Pay ${paymentSession.amount} ${paymentSession.currency}`}
    </Button>
  );
};
