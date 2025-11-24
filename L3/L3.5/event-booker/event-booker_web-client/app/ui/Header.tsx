 'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { AppBar, Toolbar, Typography, Button, Box, IconButton, Menu, MenuItem, CircularProgress } from '@mui/material'
import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import Brightness4Icon from '@mui/icons-material/Brightness4'
import Brightness7Icon from '@mui/icons-material/Brightness7'
import LogoutIcon from '@mui/icons-material/Logout'
import { useAuth } from '../context/AuthContext'
import { useThemeContext } from './ThemeProvider'

export default function Header() {
  const { user, loading, logout } = useAuth()
  const { mode, toggleTheme } = useThemeContext()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const router = useRouter()

  const open = Boolean(anchorEl)

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = async () => {
    await logout()
    router.push('/auth')
    handleMenuClose()
  }

  const handleProfile = () => {
    router.push('/profile')
    handleMenuClose()
  }

  const handleToggleTheme = () => {
    toggleTheme()
    handleMenuClose()
  }

  const handleLogin = () => {
    router.push('/auth')
  }

  return (
    <AppBar 
      position="static" 
      sx={{ 
        boxShadow: '0px 2px 4px -1px rgba(0,0,0,0.2)',
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
          <Link href="/" style={{ textDecoration: 'none', color: 'inherit' }}>
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
            >
              Event Booker
            </Typography>
          </Link>
          
          <Link href="/events" style={{ textDecoration: 'none', color: 'inherit' }}>
            <Button 
              color="inherit" 
              sx={{ textTransform: 'none' }}
            >
              Мероприятия
            </Button>
          </Link>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {loading ? (
            <CircularProgress size={24} sx={{ color: 'white' }} />
          ) : user ? (
            <>
              <Typography 
                variant="body2" 
                sx={{ 
                  color: 'white', 
                  display: { xs: 'none', sm: 'block' },
                  mr: 1
                }}
              >
                {user.name}
              </Typography>
              
              <IconButton
                size="large"
                edge="end"
                color="inherit"
                onClick={handleMenuOpen}
                aria-controls={open ? 'account-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
              >
                <AccountCircleIcon />
              </IconButton>
              
              <Menu
                id="account-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleMenuClose}
                onClick={handleMenuClose}
                PaperProps={{
                  elevation: 0,
                  sx: {
                    overflow: 'visible',
                    filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
                    mt: 1.5,
                    '& .MuiAvatar-root': {
                      width: 32,
                      height: 32,
                      ml: -0.5,
                      mr: 1,
                    },
                    '&:before': {
                      content: '""',
                      display: 'block',
                      position: 'absolute',
                      top: 0,
                      right: 14,
                      width: 10,
                      height: 10,
                      bgcolor: 'background.paper',
                      transform: 'translateY(-50%) rotate(45deg)',
                      zIndex: 0,
                    },
                  },
                }}
                transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
              >
                <MenuItem onClick={handleProfile}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <AccountCircleIcon fontSize="small" />
                    Профиль
                  </Box>
                </MenuItem>
                <MenuItem onClick={handleToggleTheme}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {mode === 'dark' ? <Brightness7Icon fontSize="small" /> : <Brightness4Icon fontSize="small" />}
                    {mode === 'dark' ? 'Светлая тема' : 'Тёмная тема'}
                  </Box>
                </MenuItem>
                <MenuItem onClick={handleLogout}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <LogoutIcon fontSize="small" />
                    Выйти
                  </Box>
                </MenuItem>
              </Menu>
            </>
          ) : (
            <Button 
              color="inherit" 
              onClick={handleLogin}
              sx={{ textTransform: 'none' }}
            >
              Войти
            </Button>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  )
}