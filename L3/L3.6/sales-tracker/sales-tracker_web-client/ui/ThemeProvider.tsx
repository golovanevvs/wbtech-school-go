"use client"

import { createContext, useContext, useState, useEffect, ReactNode } from "react"
import { ThemeProvider as MuiThemeProvider } from "@mui/material/styles"
import CssBaseline from "@mui/material/CssBaseline"
import { lightTheme, darkTheme } from "../libs/theme"

interface ThemeContextType {
  mode: "light" | "dark"
  toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

export function useThemeContext() {
  const context = useContext(ThemeContext)
  if (context === undefined) {
    throw new Error("useThemeContext must be used within a ThemeProvider")
  }
  return context
}

interface ThemeProviderProps {
  children: ReactNode
}

export default function ThemeProvider({ children }: ThemeProviderProps) {
  const [mode, setMode] = useState<"light" | "dark">("light")

  useEffect(() => {
    // Load theme from localStorage or system preference
    const savedMode = localStorage.getItem("theme-mode") as "light" | "dark" | null
    if (savedMode) {
      setMode(savedMode)
    } else {
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches
      setMode(prefersDark ? "dark" : "light")
    }
  }, [])

  const toggleTheme = () => {
    const newMode = mode === "light" ? "dark" : "light"
    setMode(newMode)
    localStorage.setItem("theme-mode", newMode)
  }

  const theme = mode === "light" ? lightTheme : darkTheme

  return (
    <ThemeContext.Provider value={{ mode, toggleTheme }}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  )
}