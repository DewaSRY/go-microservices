import { Bus, Truck, Crown } from "lucide-react";
import { Car } from "lucide-react";
import { CarPackageSlug } from "@/types/types";
import { PackageSlug } from "@/types/common";
export const PackageSlugMeta: Record<
  PackageSlug,
  {
    name: string;
    icon: React.ReactNode;
    description: string;
    slug: PackageSlug;
  }
> = {
  [PackageSlug.SEDAN]: {
    name: "Sedan",
    icon: <Car />,
    description: "Economic and comfortable",
    slug: PackageSlug.SEDAN,
  },
  [PackageSlug.SUV]: {
    name: "SUV",
    icon: <Truck />,
    description: "Spacious ride for groups",
    slug: PackageSlug.SUV,
  },
  [PackageSlug.VAN]: {
    name: "Van",
    icon: <Bus />,
    description: "Perfect for larger groups",
    slug: PackageSlug.VAN,
  },
  [PackageSlug.LUXURY]: {
    name: "Luxury",
    icon: <Crown />,
    description: "Premium experience",
    slug: PackageSlug.LUXURY,
  },
};
