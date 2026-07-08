import { Suspense } from "react";
import Header from "@/components/Header";
import Providers from "./providers";
import AuthErrorModal from "@/features/auth/components/AuthErrorModal";
import UserMenu from "@/features/auth/components/UserMenu";
import UserMenuSkeleton from "@/features/auth/components/UserMenuSkeleton";
import "react-datepicker/dist/react-datepicker.css";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <Providers>
      <Header
        userMenu={
          <Suspense fallback={<UserMenuSkeleton />}>
            <UserMenu />
          </Suspense>
        }
      />
      {children}
      <AuthErrorModal />
    </Providers>
  );
}
