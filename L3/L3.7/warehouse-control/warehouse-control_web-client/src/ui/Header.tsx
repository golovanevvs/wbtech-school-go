"use client"

import { useState } from "react"
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Menu,
  MenuItem,
  IconButton,
} from "@mui/material"
import AccountCircleIcon from "@mui/icons-material/AccountCircle"
import LoginIcon from "@mui/icons-material/Login"
import LogoutIcon from "@mui/icons-material/Logout"
import Brightness4Icon from "@mui/icons-material/Brightness4"
import Brightness7Icon from "@mui/icons-material/Brightness7"
import { useAuth } from "@/lib/contexts/AuthContext"
import { useThemeContext } from "@/lib/components/ThemeProvider"
import { useRouter } from "next/navigation"

export default function Header() {
  const { isAuthenticated, user, logout } = useAuth()
  const { mode, toggleTheme } = useThemeContext()
  const router = useRouter()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleProfileMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    logout()
    handleProfileMenuClose()
    router.push("/auth")
  }

  const handleProfile = () => {
    handleProfileMenuClose()
    router.push("/profile")
  }

  const handleLogin = () => {
    router.push("/auth")
  }

  const handleHome = () => {
    router.push("/")
  }

  const handleToggleTheme = () => {
    toggleTheme()
    handleProfileMenuClose()
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
        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
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
            Warehouse Control
          </Typography>
        </Box>

        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          {isAuthenticated && user ? (
            <>
              <IconButton
                color="inherit"
                onClick={handleProfileMenuOpen}
                aria-label="profile"
              >
                <AccountCircleIcon />
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleProfileMenuClose}
                anchorOrigin={{
                  vertical: "bottom",
                  horizontal: "right",
                }}
                transformOrigin={{
                  vertical: "top",
                  horizontal: "right",
                }}
              >
                <MenuItem onClick={handleProfile}>
                  <IconButton color="inherit" sx={{ mr: 1 }}>
                    <AccountCircleIcon />
                  </IconButton>
                  Профиль
                </MenuItem>
                <MenuItem onClick={handleToggleTheme}>
                  <IconButton color="inherit" sx={{ mr: 1 }}>
                    {mode === "dark" ? <Brightness7Icon /> : <Brightness4Icon />}
                  </IconButton>
                  {mode === "dark" ? "Светлая тема" : "Тёмная тема"}
                </MenuItem>
                <MenuItem onClick={handleLogout}>
                  <IconButton color="inherit" sx={{ mr: 1 }}>
                    <LogoutIcon />
                  </IconButton>
                  Выйти
                </MenuItem>
              </Menu>
            </>
          ) : (
            <IconButton
              color="inherit"
              onClick={handleLogin}
              aria-label="login"
            >
              <LoginIcon />
            </IconButton>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  )
}
