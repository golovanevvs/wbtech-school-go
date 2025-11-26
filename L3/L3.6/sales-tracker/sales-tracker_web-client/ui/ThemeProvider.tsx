"use client"

import { createContext, useContext, useState, useEffect, useLayoutEffect, ReactNode } from "react"
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
  // Initialize with a consistent default to prevent hydration mismatch
  const [mode, setMode] = useState<"light" | "dark">("light")
  const [mounted, setMounted] = useState(false)

  // Use useLayoutEffect for mounted state to prevent hydration mismatch
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useLayoutEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    // Load theme from localStorage or system preference
    try {
      const savedMode = localStorage.getItem("theme-mode") as "light" | "dark" | null
      if (savedMode) {
        setMode(savedMode) // eslint-disable-line react-hooks/exhaustive-deps
      } else {
        const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches
        setMode(prefersDark ? "dark" : "light") // eslint-disable-line react-hooks/exhaustive-deps
      }
    } catch (error) {
      // Fallback to light theme if localStorage is not available
      console.warn("Could not access localStorage:", error)
    }
  }, [])

  const toggleTheme = () => {
    const newMode = mode === "light" ? "dark" : "light"
    setMode(newMode)
    
    try {
      localStorage.setItem("theme-mode", newMode)
    } catch (error) {
      console.warn("Could not save theme to localStorage:", error)
    }
  }

  // Don't render theme-dependent content until mounted to prevent hydration mismatch
  if (!mounted) {
    return (
      <ThemeContext.Provider value={{ mode: "light", toggleTheme: () => {} }}>
        <MuiThemeProvider theme={lightTheme}>
          <CssBaseline />
          {children}
        </MuiThemeProvider>
      </ThemeContext.Provider>
    )
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