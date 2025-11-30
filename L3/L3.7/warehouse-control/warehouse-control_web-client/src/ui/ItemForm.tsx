"use client"

import {
  Box,
  Paper,
  TextField,
  Button,
  Typography,
  CircularProgress,
} from "@mui/material"
import { Item } from "@/lib/types/items"
import { useState } from "react"

interface ItemFormProps {
  initialData?: Partial<Item>
  onSubmit: (
    data: Omit<Item, "id" | "created_at" | "updated_at">
  ) => Promise<void>
  onCancel: () => void
  title: string
  submitButtonText: string
  loading?: boolean
}

export default function ItemForm({
  initialData,
  onSubmit,
  onCancel,
  title,
  submitButtonText,
  loading = false,
}: ItemFormProps) {
  const [formData, setFormData] = useState({
    name: initialData?.name || "",
    price: initialData?.price || 0,
    quantity: initialData?.quantity || 0,
  })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = "Название товара обязательно"
    }

    if (formData.price < 0) {
      newErrors.price = "Цена не может быть отрицательной"
    }

    if (formData.quantity < 0) {
      newErrors.quantity = "Количество не может быть отрицательным"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      setSubmitting(true)
      await onSubmit(formData)
    } catch (error) {
      console.error("Form submission error:", error)
    } finally {
      setSubmitting(false)
    }
  }

  const handleInputChange = (field: string, value: string | number) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }))

    if (errors[field]) {
      setErrors((prev) => ({
        ...prev,
        [field]: "",
      }))
    }
  }

  return (
    <Box sx={{ maxWidth: 600, mx: "auto", p: 2 }}>
      <Paper elevation={2} sx={{ p: 4 }}>
        <Typography variant="h5" sx={{ mb: 3, textAlign: "center" }}>
          {title}
        </Typography>

        <Box component="form" onSubmit={handleSubmit}>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
            <TextField
              fullWidth
              label="Название товара"
              value={formData.name}
              onChange={(e) => handleInputChange("name", e.target.value)}
              error={!!errors.name}
              helperText={errors.name}
              disabled={submitting || loading}
              required
            />

            <Box sx={{ display: "flex", gap: 2 }}>
              <TextField
                fullWidth
                label="Цена (руб.)"
                type="text"
                value={formData.price}
                onChange={(e) => {
                  const value = e.target.value
                  const cleanValue = value
                    .replace(/[^0-9.]/g, "")
                    .replace(/(\..*)\./g, "$1")
                  const numValue =
                    cleanValue === "" ? 0 : parseFloat(cleanValue)
                  handleInputChange("price", numValue)
                }}
                error={!!errors.price}
                helperText={errors.price || "Введите цену (например: 123.45)"}
                disabled={submitting || loading}
                required
              />

              <TextField
                fullWidth
                label="Количество"
                type="text"
                value={formData.quantity}
                onChange={(e) => {
                  const value = e.target.value
                  const cleanValue = value.replace(/[^0-9]/g, "")
                  const numValue = cleanValue === "" ? 0 : parseInt(cleanValue)
                  handleInputChange("quantity", numValue)
                }}
                error={!!errors.quantity}
                helperText={errors.quantity || "Введите количество"}
                disabled={submitting || loading}
                required
              />
            </Box>

            <Box
              sx={{
                display: "flex",
                gap: 2,
                justifyContent: { xs: "center", md: "flex-end" },
                mt: 2,
                flexWrap: { xs: "wrap", md: "nowrap" },
              }}
            >
              <Button
                variant="outlined"
                onClick={onCancel}
                disabled={submitting || loading}
                sx={{
                  minWidth: { xs: "45%", md: "auto" },
                  flex: { xs: 1, md: "none" },
                }}
              >
                Отменить
              </Button>
              <Button
                type="submit"
                variant="contained"
                disabled={submitting || loading}
                startIcon={submitting ? <CircularProgress size={20} /> : null}
                sx={{
                  minWidth: { xs: "45%", md: "auto" },
                  flex: { xs: 1, md: "none" },
                }}
              >
                {submitting ? "Сохранение..." : submitButtonText}
              </Button>
            </Box>
          </Box>
        </Box>
      </Paper>
    </Box>
  )
}
