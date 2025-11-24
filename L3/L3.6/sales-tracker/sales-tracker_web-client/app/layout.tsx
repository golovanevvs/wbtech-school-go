import type { Metadata } from "next";
import { Geist } from "next/font/google";
import "./globals.css";
import ThemeProvider from "../ui/ThemeProvider";
import Header from "../ui/Header";
import Footer from "../ui/Footer";
import { Box, Container } from "@mui/material";

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
});

export const metadata: Metadata = {
  title: "Sales Tracker",
  description: "Система отслеживания продаж с аналитикой",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru">
      <body className={`${geist.variable}`}>
        <ThemeProvider>
          <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <Header />
            <Container maxWidth="lg" sx={{ flex: 1, py: 2 }}>
              {children}
            </Container>
            <Footer />
          </Box>
        </ThemeProvider>
      </body>
    </html>
  );
}
