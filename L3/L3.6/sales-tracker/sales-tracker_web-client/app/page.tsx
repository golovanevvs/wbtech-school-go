"use client"

import { Box } from "@mui/material"
import { useRouter } from "next/navigation"
import Button from "../ui/Button"

export default function Home() {
  const router = useRouter()

  const handleNavigate = (path: string) => {
    router.push(path)
  }

  const menuItems = [
    {
      title: "Добавить запись",
      path: "/add-item",
      color: "primary" as const,
    },
    {
      title: "Посмотреть все записи",
      path: "/items",
      color: "secondary" as const,
    },
    {
      title: "Получить аналитику",
      path: "/analytics",
      color: "success" as const,
    },
  ]

  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center',
      gap: 3
    }}>
      {menuItems.map((item, index) => (
        <Button
          key={index}
          variant="contained"
          color={item.color}
          onClick={() => handleNavigate(item.path)}
          sx={{ 
            width: '100%',
            maxWidth: '300px',
            minHeight: '56px',
            fontSize: '1.1rem'
          }}
        >
          {item.title}
        </Button>
      ))}
    </Box>
  )
}