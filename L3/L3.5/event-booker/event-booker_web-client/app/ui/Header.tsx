 'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { AppBar, Toolbar, Typography, Button, Box, IconButton, Menu, MenuItem, CircularProgress } from '@mui/material'
import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import { useAuth } from '../context/AuthContext'

export default function Header() {
  const { user, loading, logout } = useAuth()
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

  // 👇 Добавьте эту функцию для перехода на страницу логина
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
            onClick={() => router.push('/')}
          >
            Event Booker
          </Typography>
          
          <Button 
            color="inherit" 
            onClick={() => router.push('/events')}
            sx={{ textTransform: 'none' }}
          >
            Мероприятия
          </Button>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {loading ? (
            // Показываем индикатор загрузки при проверке статуса аутентификации
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
                  Профиль
                </MenuItem>
                <MenuItem onClick={handleLogout}>
                  Выйти
                </MenuItem>
              </Menu>
            </>
          ) : (
            <Button 
              color="inherit" 
              onClick={handleLogin} // 👈 Используем новую функцию
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