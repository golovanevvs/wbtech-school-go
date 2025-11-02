import type { Metadata } from "next"
import { Geist, Geist_Mono } from "next/font/google"
// import "./globals.css";
import { AppRouterCacheProvider } from "@mui/material-nextjs/v15-appRouter"
import ThemeProvider from "@/app/ui/ThemeProvider"
import styles from "./page.module.css"
import ThemeToggle from "./ui/ThemeToggle"
import Footer from "./ui/Footer"

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
})

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
})

export const metadata: Metadata = {
  title: "Delayed notifier",
  description: "Delayed notifier",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="ru">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </head>
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        <AppRouterCacheProvider>
          <ThemeProvider>
            <ThemeToggle />
            <main className={styles.page}>{children}</main>
            <Footer />
          </ThemeProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  )
}
