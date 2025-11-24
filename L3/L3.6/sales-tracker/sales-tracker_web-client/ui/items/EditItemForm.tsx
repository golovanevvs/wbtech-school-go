"use client"

import { useState, useEffect } from "react"
import { Box, Typography, Alert, Dialog, DialogTitle, DialogContent, DialogActions, MenuItem } from "@mui/material"
import Button from "../Button"
import Input from "../Input"
import { SalesRecord, SalesRecordFormData } from "../../libs/types"

interface EditItemFormProps {
  record: SalesRecord
  onSave: (id: number, data: Partial<SalesRecord>) => void
  onCancel: () => void
}

export default function EditItemForm({ record, onSave, onCancel }: EditItemFormProps) {
  const [formData, setFormData] = useState<SalesRecordFormData>({
    type: record.type,
    category: record.category,
    date: new Date(record.date).toISOString().slice(0, 16), // Format for datetime-local input
    amount: record.amount,
  })
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

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

  useEffect(() => {
    setFormData({
      type: record.type,
      category: record.category,
      date: new Date(record.date).toISOString().slice(0, 16),
      amount: record.amount,
    })
  }, [record])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
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

      await onSave(record.id, formData)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка при обновлении записи")
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
    <Dialog open={true} onClose={onCancel} maxWidth="sm" fullWidth>
      <DialogTitle>
        Редактировать запись #{record.id}
      </DialogTitle>
      
      <Box component="form" onSubmit={handleSubmit}>
        <DialogContent>
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
            <Alert severity="error" sx={{ mt: 2 }}>
              {error}
            </Alert>
          )}
        </DialogContent>

        <DialogActions>
          <Button onClick={onCancel} variant="outlined">
            Отмена
          </Button>
          <Button type="submit" disabled={loading} variant="contained">
            {loading ? "Сохранение..." : "Сохранить"}
          </Button>
        </DialogActions>
      </Box>
    </Dialog>
  )
}