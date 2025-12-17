import { PropsWithChildren } from "react";

export default function TitleBadge({ children }: PropsWithChildren) {
  return (
    <div className="w-full">
      <div className=" w-11/12 m-auto py-4 flex justify-center items-center  bg-gray-600/90 text-white rounded-sm">
        {children}
      </div>
    </div>
  );
}
