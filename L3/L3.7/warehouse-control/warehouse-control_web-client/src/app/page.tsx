"use client"

import { useState, useEffect } from "react"
import { useAuth } from "@/lib/contexts/AuthContext"
import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { itemsAPI } from "@/lib/api/items"
import { Item } from "@/lib/types/items"
import { 
  Box, 
  Typography, 
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Chip,
  CircularProgress,
  Alert
} from "@mui/material"
import EditIcon from "@mui/icons-material/Edit"
import DeleteIcon from "@mui/icons-material/Delete"
import HistoryIcon from "@mui/icons-material/History"

export default function Home() {
  const { user, hasRole } = useAuth()
  const { isLoading, isAuthenticated } = useAuthGuard()
  
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Загрузка товаров
  useEffect(() => {
    const loadItems = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await itemsAPI.getItems()
        setItems(response.items || [])
      } catch (err) {
        console.error("Failed to load items:", err)
        setError(err instanceof Error ? err.message : "Не удалось загрузить товары")
      } finally {
        setLoading(false)
      }
    }

    if (isAuthenticated && !isLoading) {
      loadItems()
    }
  }, [isAuthenticated, isLoading])

  // Форматирование даты
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU')
  }

  // Форматирование цены
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB'
    }).format(price)
  }

  // Обработчики действий
  const handleEdit = (itemId: number) => {
    console.log("Edit item:", itemId)
    // TODO: Реализовать редактирование
  }

  const handleDelete = (itemId: number) => {
    console.log("Delete item:", itemId)
    // TODO: Реализовать удаление
  }

  const handleHistory = (itemId: number) => {
    console.log("View history for item:", itemId)
    // TODO: Реализовать просмотр истории
  }

  // Показываем загрузку
  if (isLoading || loading) {
    return (
      <Box 
        sx={{ 
          display: "flex", 
          justifyContent: "center", 
          alignItems: "center", 
          minHeight: "50vh" 
        }}
      >
        <CircularProgress />
      </Box>
    )
  }

  // Если не авторизован, useAuthGuard перенаправит на /auth
  if (!isAuthenticated) {
    return null
  }

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto", p: 3 }}>
      <Typography variant="h4" sx={{ mb: 3, textAlign: "center" }}>
        Список товаров склада
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Информация о пользователе и роль - НАД ТАБЛИЦЕЙ */}
      <Paper 
        elevation={1} 
        sx={{ 
          p: 3, 
          mb: 3, 
          backgroundColor: "background.paper",
          border: "1px solid",
          borderColor: "divider"
        }}
      >
        <Typography variant="h6" sx={{ mb: 1 }}>
          Здравствуйте, {user?.name}!
        </Typography>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, flexWrap: "wrap" }}>
          <Typography variant="body1">
            Ваша роль:
          </Typography>
          <Chip 
            label={user?.user_role} 
            color="primary" 
            size="medium"
            sx={{ fontWeight: "medium" }}
          />
          <Typography variant="body2" color="text.secondary">
            <strong>Ваши права доступа:</strong>
            {hasRole(["Кладовщик"]) && " • Редактирование и удаление товаров"}
            {hasRole(["Менеджер"]) && " • Просмотр товаров"}
            {hasRole(["Аудитор"]) && " • Просмотр истории изменений"}
          </Typography>
        </Box>
      </Paper>

      <TableContainer 
        component={Paper} 
        elevation={2}
        sx={{ 
          backgroundColor: "background.paper",
          border: "1px solid",
          borderColor: "divider"
        }}
      >
        <Table>
          <TableHead>
            <TableRow sx={{ backgroundColor: "action.hover" }}>
              <TableCell sx={{ fontWeight: "bold" }}>ID</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Название</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Цена</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Количество</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Дата создания</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Дата изменения</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} sx={{ textAlign: "center", py: 4 }}>
                  <Typography color="text.secondary">
                    Товары не найдены
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              items.map((item) => (
                <TableRow key={item.id} hover sx={{ "&:hover": { backgroundColor: "action.hover" } }}>
                  <TableCell>{item.id}</TableCell>
                  <TableCell sx={{ fontWeight: "medium" }}>{item.name}</TableCell>
                  <TableCell>{formatPrice(item.price)}</TableCell>
                  <TableCell>
                    <Chip 
                      label={item.quantity.toString()} 
                      color={item.quantity > 0 ? "success" : "error"}
                      size="small"
                      sx={{ 
                        color: item.quantity > 0 ? "success.contrastText" : "error.contrastText",
                        "& .MuiChip-label": {
                          color: item.quantity > 0 ? "success.contrastText" : "error.contrastText",
                        }
                      }}
                    />
                  </TableCell>
                  <TableCell>{formatDate(item.created_at)}</TableCell>
                  <TableCell>{formatDate(item.updated_at)}</TableCell>
                  <TableCell>
                    <Box sx={{ display: "flex", gap: 0.5 }}>
                      {/* Кнопки для Кладовщика */}
                      {hasRole(["Кладовщик"]) && (
                        <>
                          <IconButton 
                            size="small" 
                            color="primary"
                            onClick={() => handleEdit(item.id)}
                            title="Редактировать"
                            sx={{
                              "&:hover": {
                                backgroundColor: "primary.main",
                                color: "primary.contrastText"
                              }
                            }}
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                          <IconButton 
                            size="small" 
                            color="error"
                            onClick={() => handleDelete(item.id)}
                            title="Удалить"
                            sx={{
                              "&:hover": {
                                backgroundColor: "error.main",
                                color: "error.contrastText"
                              }
                            }}
                          >
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </>
                      )}
                      
                      {/* Кнопка истории для Аудитора */}
                      {hasRole(["Аудитор"]) && (
                        <IconButton 
                          size="small" 
                          color="info"
                          onClick={() => handleHistory(item.id)}
                          title="История изменений"
                          sx={{
                            "&:hover": {
                              backgroundColor: "info.main",
                              color: "info.contrastText"
                            }
                          }}
                        >
                          <HistoryIcon fontSize="small" />
                        </IconButton>
                      )}
                    </Box>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  )
}
