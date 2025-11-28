"use client"

import { useState } from "react"
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Button,
  Menu,
  MenuItem,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
} from "@mui/material"
import MenuIcon from "@mui/icons-material/Menu"
import { useAuth } from "@/lib/contexts/AuthContext"
import ThemeToggle from "./ThemeToggle"
import { useRouter } from "next/navigation"

export default function Header() {
  const { isAuthenticated, user, logout, hasRole } = useAuth()
  const router = useRouter()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [menuDrawerOpen, setMenuDrawerOpen] = useState(false)

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuOpen = () => {
    setMenuDrawerOpen(true)
  }

  const handleMenuClose = () => {
    setMenuDrawerOpen(false)
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

  const handleItems = () => {
    handleMenuClose()
    router.push("/items")
  }

  const handleHistory = () => {
    handleMenuClose()
    router.push("/history")
  }

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

          {isAuthenticated && (
            <IconButton
              color="inherit"
              onClick={handleMenuOpen}
              sx={{ ml: -1 }}
            >
              <MenuIcon />
            </IconButton>
          )}
        </Box>

        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          {isAuthenticated && user ? (
            <>
              <Button
                color="inherit"
                onClick={handleProfileMenuOpen}
                sx={{ textTransform: "none" }}
              >
                Профиль
              </Button>
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
                <MenuItem onClick={handleProfile}>Профиль</MenuItem>
                <MenuItem onClick={handleProfileMenuClose}>
                  <ThemeToggle />
                </MenuItem>
                <MenuItem onClick={handleLogout}>Выйти</MenuItem>
              </Menu>
            </>
          ) : (
            <ThemeToggle />
          )}
        </Box>
      </Toolbar>

      {/* Выдвижное меню слева */}
      <Drawer
        anchor="left"
        open={menuDrawerOpen}
        onClose={handleMenuClose}
        PaperProps={{
          sx: {
            width: 250,
            mt: 8, // Отступ сверху для учета AppBar
          },
        }}
      >
        <List>
          {hasRole(["Кладовщик", "Менеджер"]) && (
            <ListItem disablePadding>
              <ListItemButton onClick={handleItems}>
                <ListItemText primary="Список товаров" />
              </ListItemButton>
            </ListItem>
          )}
          <ListItem disablePadding>
            <ListItemButton onClick={handleHistory}>
              <ListItemText primary="История действий" />
            </ListItemButton>
          </ListItem>
        </List>
      </Drawer>
    </AppBar>
  )
}
