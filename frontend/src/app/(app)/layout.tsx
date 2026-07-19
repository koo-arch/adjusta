import { Suspense } from "react";
import Header from "@/components/layout/Header";
import HeaderSkeleton from "@/components/layout/Header/HeaderSkeleton";
import Providers from "./providers";
import AuthErrorModal from "@/features/auth/components/AuthErrorModal";
import UserMenu from "@/features/auth/components/UserMenu";
import UserMenuSkeleton from "@/features/auth/components/UserMenuSkeleton";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <Providers>
      {/* Header は usePathname を使うため、動的ルート([id])のプリレンダーでは
          suspend してこの fallback が静的シェルに入る。静的ルートでは実物が入る */}
      <Suspense fallback={<HeaderSkeleton />}>
        <Header
          userMenu={
            <Suspense fallback={<UserMenuSkeleton />}>
              <UserMenu />
            </Suspense>
          }
        />
      </Suspense>
      {children}
      <AuthErrorModal />
    </Providers>
  );
}
