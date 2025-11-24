"use client"

import { useState } from "react"
import { Box, Typography, Alert, MenuItem } from "@mui/material"
import Card from "../Card"
import Button from "../Button"
import Input from "../Input"
import { SalesRecordFormData } from "../../libs/types"
import { createSalesRecord } from "../../libs/api/sales"

interface AddItemFormProps {
  onSuccess?: () => void
}

export default function AddItemForm({ onSuccess }: AddItemFormProps) {
  const [formData, setFormData] = useState<SalesRecordFormData>({
    type: "income",
    category: "",
    date: new Date().toISOString().slice(0, 16), // ISO format for datetime-local input
    amount: 0,
  })
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const categories = [
    "Продажи",
    "Услуги",
    "Инвестиции",
    "Еда",
    "Транспорт",
    "Жилье",
    "Развлечения",
    "Здоровье",
    "Образование",
    "Другое",
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setSuccess(false)
    setLoading(true)

    try {
      // Validate form data
      if (!formData.category.trim()) {
        setError("Категория обязательна для заполнения")
        setLoading(false)
        return
      }

      if (formData.amount <= 0) {
        setError("Сумма должна быть положительным числом")
        setLoading(false)
        return
      }

      await createSalesRecord(formData)
      
      setSuccess(true)
      
      // Reset form
      setFormData({
        type: "income",
        category: "",
        date: new Date().toISOString().slice(0, 16),
        amount: 0,
      })

      if (onSuccess) {
        onSuccess()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка при добавлении записи")
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field: keyof SalesRecordFormData) => (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const value = field === "amount" ? parseFloat(e.target.value) || 0 : e.target.value
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  return (
    <Card>
      <Typography variant="h5" component="h1" gutterBottom>
        Добавить запись
      </Typography>

      <Box component="form" onSubmit={handleSubmit} sx={{ width: "100%" }}>
        <Input
          select
          label="Тип"
          value={formData.type}
          onChange={handleChange("type")}
          required
        >
          <MenuItem value="income">Доход</MenuItem>
          <MenuItem value="expense">Расход</MenuItem>
        </Input>

        <Input
          select
          label="Категория"
          value={formData.category}
          onChange={handleChange("category")}
          required
        >
          {categories.map((category) => (
            <MenuItem key={category} value={category}>
              {category}
            </MenuItem>
          ))}
        </Input>

        <Input
          label="Дата"
          type="datetime-local"
          value={formData.date}
          onChange={handleChange("date")}
          required
          InputLabelProps={{
            shrink: true,
          }}
        />

        <Input
          label="Сумма"
          type="number"
          value={formData.amount || ""}
          onChange={handleChange("amount")}
          required
          inputProps={{
            step: "0.01",
            min: "0",
          }}
        />

        {error && (
          <Alert severity="error" sx={{ mt: 2, width: "100%" }}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mt: 2, width: "100%" }}>
            Запись успешно добавлена!
          </Alert>
        )}

        <Button type="submit" disabled={loading} sx={{ mt: 2 }}>
          {loading ? "Добавление..." : "Добавить запись"}
        </Button>
      </Box>
    </Card>
  )
}