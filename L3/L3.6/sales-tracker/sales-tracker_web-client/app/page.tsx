"use client"

import { Box, Button } from "@mui/material"
import { useRouter } from "next/navigation"
import ItemsList from "../ui/items/ItemsList"

export default function Home() {
  const router = useRouter()

  const handleNavigate = (path: string) => {
    router.push(path)
  }

  return (
    <Box sx={{ width: "100%" }}>
      {/* Header with action buttons */}
      <Box sx={{ 
        display: "flex", 
        gap: 2, 
        mb: 3, 
        justifyContent: "center",
        flexWrap: "wrap"
      }}>
        <Button
          variant="contained"
          color="primary"
          onClick={() => handleNavigate("/add-item")}
          sx={{ 
            minWidth: "150px",
            fontSize: "1rem"
          }}
        >
          Добавить запись
        </Button>
        
        <Button
          variant="outlined"
          color="success"
          onClick={() => handleNavigate("/analytics")}
          sx={{ 
            minWidth: "150px",
            fontSize: "1rem"
          }}
        >
          Аналитика
        </Button>
      </Box>

      {/* Main content - Items List */}
      <ItemsList />
    </Box>
  )
}