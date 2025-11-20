import { Bus, Truck, Crown } from "lucide-react";
import { Car } from "lucide-react";
import { PackageSlug } from "@/types/common";
import { cn } from "@/libs/utils";

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

interface PackageSlugPros {
  onSelectSlug: (_data: PackageSlug) => void;
}

export default function PackageSlugWidget({ onSelectSlug }: PackageSlugPros) {
  return (
    <div className="space-y-3 sm:space-y-4">
      {Object.entries(PackageSlugMeta).map(([slug, meta]) => (
        <div
          key={slug}
          className={cn(
            "flex items-center gap-3 sm:gap-4 p-3 sm:p-4 sm:rounded-lg sm:border transition-all cursor-pointer",
            "hover:border-primary hover:bg-primary/5"
          )}
          onClick={onSelectSlug.bind(null, meta.slug)}
        >
          <div className="p-1.5 sm:p-2 bg-gray-100 rounded-lg">
            {meta?.icon}
          </div>
          <div>
            <h3 className="font-medium text-sm sm:text-base">{meta?.name}</h3>
            <p className="text-xs sm:text-sm text-gray-500">
              {meta?.description}
            </p>
          </div>
        </div>
      ))}
    </div>
  );
}
