import { useTranslations } from "next-intl";

export function HomeNavigationWidget() {
  const navigationT = useTranslations("navigation");
  return (
    <header className="xl:w-[1200px] bg-white px-4 py-1 rounded-md">
      <div>{navigationT("logo")}</div>
    </header>
  );
}
