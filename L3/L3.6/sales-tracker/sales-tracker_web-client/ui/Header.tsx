'use client'

import Link from 'next/link'
import { AppBar, Toolbar, Typography, Box } from '@mui/material'
import ThemeToggle from './ThemeToggle'

export default function Header() {
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
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
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
              Sales Tracker
            </Typography>
          </Link>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <ThemeToggle />
        </Box>
      </Toolbar>
    </AppBar>
  )
}