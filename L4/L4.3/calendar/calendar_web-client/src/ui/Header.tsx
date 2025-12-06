"use client"

import { AppBar, Toolbar, Typography, IconButton } from "@mui/material"
import Brightness4Icon from "@mui/icons-material/Brightness4"
import Brightness7Icon from "@mui/icons-material/Brightness7"
import { useThemeContext } from "@/lib/components/ThemeProvider"
import { useRouter } from "next/navigation"

export default function Header() {
  const { mode, toggleTheme } = useThemeContext()
  const router = useRouter()

  const handleHome = () => {
    router.push("/")
  }

  return (
    <AppBar
      position="static"
      sx={{
        boxShadow: "0px 4px 6px -1px rgba(0,0,0,0.2)",
        borderRadius: 0,
      }}
    >
      <Toolbar
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          px: { xs: 1, sm: 2 },
        }}
      >
        {/* Левая часть - название */}
        <Typography
          variant="h6"
          component="div"
          sx={{
            fontWeight: "bold",
            cursor: "pointer",
            "&:hover": {
              opacity: 0.8,
            },
          }}
          onClick={handleHome}
        >
          Calendar
        </Typography>

        {/* Правая часть - переключатель темы */}
        <IconButton
          color="inherit"
          onClick={toggleTheme}
          aria-label="toggle theme"
        >
          {mode === "dark" ? <Brightness7Icon /> : <Brightness4Icon />}
        </IconButton>
      </Toolbar>
    </AppBar>
  )
}
