// src/app/history/page.tsx - для всех авторизованных пользователей
"use client"

import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { Box, Typography, Card, CardContent } from "@mui/material"

export default function HistoryPage() {
  // Проверяем авторизацию для всех ролей
  useAuthGuard()

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto", mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        История действий
      </Typography>

      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Последние действия в системе
          </Typography>

          {/* Здесь будет история действий пользователей */}
          <div>
            <Typography color="text.secondary">
              История действий будет отображаться здесь
            </Typography>
          </div>
        </CardContent>
      </Card>
    </Box>
  )
}
