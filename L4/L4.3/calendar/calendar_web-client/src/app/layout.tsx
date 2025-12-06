import type { Metadata } from "next"
import { Geist, Geist_Mono } from "next/font/google"
import { AppRouterCacheProvider } from "@mui/material-nextjs/v16-appRouter"
import { Box, Stack } from "@mui/material"
import ThemeProvider from "@/lib/components/ThemeProvider"
import Header from "@/ui/Header"
import Footer from "@/ui/Footer"

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
})

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
})

export const metadata: Metadata = {
  title: "Calendar",
  description: "Calendar",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="ru">
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        <AppRouterCacheProvider>
          <ThemeProvider>
              <Box
                sx={{
                  width: "100%",
                  minHeight: "100vh",
                  px: { xs: 0, sm: 2 },
                  py: 2,
                  bgcolor: "background.default",
                  maxWidth: "100vw",
                  mx: "auto",
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <Stack
                  spacing={2}
                  alignItems="center"
                  sx={{
                    flex: 1,
                    width: "100%",
                    display: "flex",
                    flexDirection: "column",
                  }}
                >
                  <Header />
                  <main style={{ flex: 1, width: "100%" }}>{children}</main>
                  <Footer />
                </Stack>
              </Box>
          </ThemeProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  )
}
