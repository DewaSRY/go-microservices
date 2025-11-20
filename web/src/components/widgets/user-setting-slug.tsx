import PackageSlugWidget from "@components/widgets/package-slug-widget";

import useUserRideProfile from "@/hooks/state/useUserRideProfile";

export default function UserSettingSlug() {
  const { setPackageSlug } = useUserRideProfile();
  return (
    <div>
      <div>
        <PackageSlugWidget onSelectSlug={setPackageSlug} />
      </div>
    </div>
  );
}
