import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AppRouterCacheProvider } from "@mui/material-nextjs/v15-appRouter"
import ThemeProvider from "@/app/ui/ThemeProvider"
import Header from "@/app/ui/Header"
import styles from "./page.module.css"

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Delayed notifier",
  description: "Delayed notifier",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru">
          <body className={`${geistSans.variable} ${geistMono.variable}`}>
            <AppRouterCacheProvider>
              <ThemeProvider>
                <Header />
                <main className={styles.page}>{children}</main>
              </ThemeProvider>
            </AppRouterCacheProvider>
          </body>
        </html>
  );
}
