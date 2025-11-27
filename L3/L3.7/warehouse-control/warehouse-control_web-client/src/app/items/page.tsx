// src/app/items/page.tsx - только для Кладовщик и Менеджер
"use client"

import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { Box, Typography } from "@mui/material"

export default function ItemsPage() {
  // Только для ролей "Кладовщик" и "Менеджер"
  useAuthGuard(["Кладовщик", "Менеджер"])

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Список товаров
      </Typography>
      
      {/* Здесь будет список товаров */}
      <div>Здесь будет таблица товаров склада</div>
    </Box>
  )
}