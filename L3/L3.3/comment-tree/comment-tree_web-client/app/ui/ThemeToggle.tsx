"use client"

import { useThemeContext } from "@/app/ui/ThemeProvider"
import { IconButton } from "@mui/material"
import Brightness4Icon from "@mui/icons-material/Brightness4"
import Brightness7Icon from "@mui/icons-material/Brightness7"

export default function ThemeToggle() {
  const { mode, toggleTheme } = useThemeContext()

  return (
    <IconButton onClick={toggleTheme} color="inherit">
      {mode === "dark" ? <Brightness7Icon /> : <Brightness4Icon />}
    </IconButton>
  )
}
