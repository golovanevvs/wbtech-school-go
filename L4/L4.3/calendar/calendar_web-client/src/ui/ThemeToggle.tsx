"use client"

import { IconButton } from "@mui/material"
import Brightness4Icon from "@mui/icons-material/Brightness4"
import Brightness7Icon from "@mui/icons-material/Brightness7"
import { useThemeContext } from "../lib/components/ThemeProvider"

export default function ThemeToggle() {
  const { mode, toggleTheme } = useThemeContext()

  return (
    <IconButton
      color="inherit"
      onClick={toggleTheme}
      sx={{ ml: 1 }}
      aria-label="toggle theme"
    >
      {mode === "dark" ? <Brightness7Icon /> : <Brightness4Icon />}
    </IconButton>
  )
}
