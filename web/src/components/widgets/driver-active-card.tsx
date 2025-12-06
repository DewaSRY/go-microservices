import { HTMLAttributes } from "react";

export interface DriverActiveCardProps extends HTMLAttributes<HTMLDivElement> {
  packageSlug: string;
  title: string;
}

export default function DriverActiveCard({
  packageSlug,
  title,
}: DriverActiveCardProps) {
  return (
    <div className="mb-2">
      <div>
        <h3 className="font-medium text-gray-600">{title}</h3>
        <span className="font-extralight text-gray-500">{packageSlug}</span>
      </div>
    </div>
  );
}
