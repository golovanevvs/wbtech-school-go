import type { Metadata } from "next"
import { Geist, Geist_Mono } from "next/font/google"
import { AppRouterCacheProvider } from "@mui/material-nextjs/v15-appRouter"
import ThemeProvider from "./ui/ThemeProvider"
import ThemeToggle from "./ui/ThemeToggle"
import Header from "./ui/Header"
import Footer from "./ui/Footer"
import { Box, Stack } from "@mui/material"

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
})

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
})

export const metadata: Metadata = {
  title: "Event Booker",
  description: "Event Booker",
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
            <Box
              sx={{
                width: "100%",
                minHeight: "100vh",
                px: { xs: 0, sm: 2 },
                py: 2,
                bgcolor: "background.default",
                maxWidth: "100vw",
                mx: "auto",
              }}
            >
              <Stack spacing={2} alignItems="center">
                <Header />
                <main>{children}</main>
                <Footer />
              </Stack>
            </Box>
          </ThemeProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  )
}
