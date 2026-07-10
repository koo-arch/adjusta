import type { Metadata } from "next";
import ThemeProvider from "@/components/ThemeProvider";
import { Inter, Noto_Sans_JP } from "next/font/google";
import "./globals.css";

// 欧文は Inter、和文は Noto Sans JP が受ける混植構成(DESIGN.md 3.3)
const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

const notoSansJP = Noto_Sans_JP({
  subsets: ["latin"],
  weight: ["400", "500", "700"],
  variable: "--font-noto-sans-jp",
});

export const metadata: Metadata = {
  title: "Adjusta",
  description: "イベントの日程調整をもっとシンプルにする Adjusta",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="ja"
      suppressHydrationWarning
    >
      <body className={`${inter.variable} ${notoSansJP.variable} font-sans`}>
        <ThemeProvider
          defaultTheme="system"
          enableSystem
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
