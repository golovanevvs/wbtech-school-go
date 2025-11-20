"use client"

import { createContext, useContext, useState, useMemo } from "react"
import { ThemeProvider as MuiThemeProvider } from "@mui/material/styles"
import CssBaseline from "@mui/material/CssBaseline"
import { PaletteMode } from "@mui/material"
import { lightTheme, darkTheme } from "@/app/lib/theme"

type ThemeContextType = {
  mode: PaletteMode
  toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

export function useThemeContext() {
  const context = useContext(ThemeContext)
  if (!context) {
    throw new Error("useThemeContext must be used within a ThemeProvider")
  }
  return context
}

export default function ThemeProvider({
  children,
}: {
  children: React.ReactNode
}) {
  const [mode, setMode] = useState<PaletteMode>("light")

  const toggleTheme = () => {
    setMode((prev) => (prev === "light" ? "dark" : "light"))
  }

  const theme = useMemo(
    () => (mode === "light" ? lightTheme : darkTheme),
    [mode]
  )

  return (
    <ThemeContext.Provider value={{ mode, toggleTheme }}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  )
}