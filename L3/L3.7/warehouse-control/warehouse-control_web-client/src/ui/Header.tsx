'use client'

import Link from 'next/link'
import { useState } from 'react'
import { AppBar, Toolbar, Typography, Box, Button, Menu, MenuItem } from '@mui/material'
import { useAuth } from '@/lib/contexts/AuthContext'
import ThemeToggle from './ThemeToggle'
import { useRouter } from 'next/navigation'

export default function Header() {
  const { isAuthenticated, user, logout, hasRole } = useAuth()
  const router = useRouter()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null)

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setMenuAnchorEl(null)
  }

  const handleProfileMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    logout()
    handleProfileMenuClose()
    router.push('/auth')
  }

  const handleProfile = () => {
    handleProfileMenuClose()
    router.push('/profile')
  }

  const handleItems = () => {
    handleMenuClose()
    router.push('/items')
  }

  const handleHistory = () => {
    handleMenuClose()
    router.push('/history')
  }

  const handleHome = () => {
    router.push('/')
  }

  return (
    <AppBar 
      position="static" 
      sx={{ 
        boxShadow: '0px 4px 6px -1px rgba(0,0,0,0.2)',
        borderRadius: 0,
      }}
    >
      <Toolbar 
        sx={{ 
          display: 'flex', 
          justifyContent: 'space-between',
          alignItems: 'center',
          px: { xs: 1, sm: 2 }
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Typography 
            variant="h6" 
            component="div" 
            sx={{ 
              fontWeight: 'bold',
              cursor: 'pointer',
              '&:hover': {
                opacity: 0.8
              }
            }}
            onClick={handleHome}
          >
            Warehouse Control
          </Typography>

          {isAuthenticated && (
            <>
              <Button 
                color="inherit" 
                onClick={handleMenuOpen}
                sx={{ textTransform: 'none' }}
              >
                Меню
              </Button>
              <Menu
                anchorEl={menuAnchorEl}
                open={Boolean(menuAnchorEl)}
                onClose={handleMenuClose}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'left',
                }}
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'left',
                }}
              >
                {hasRole(['Кладовщик', 'Менеджер']) && (
                  <MenuItem onClick={handleItems}>
                    Список товаров
                  </MenuItem>
                )}
                <MenuItem onClick={handleHistory}>
                  История действий
                </MenuItem>
              </Menu>
            </>
          )}
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {isAuthenticated && user ? (
            <>
              <Button 
                color="inherit" 
                onClick={handleProfileMenuOpen}
                sx={{ textTransform: 'none' }}
              >
                Профиль
              </Button>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleProfileMenuClose}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
              >
                <MenuItem onClick={handleProfile}>
                  Профиль
                </MenuItem>
                <MenuItem onClick={handleProfileMenuClose}>
                  <ThemeToggle />
                </MenuItem>
                <MenuItem onClick={handleLogout}>
                  Выйти
                </MenuItem>
              </Menu>
            </>
          ) : (
            <ThemeToggle />
          )}
        </Box>
      </Toolbar>
    </AppBar>
  )
}