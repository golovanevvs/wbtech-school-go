"use client"

import { useState } from "react"
import { Box, Typography, Alert, MenuItem } from "@mui/material"
import Link from "next/link"
import { DatePicker } from "@mui/x-date-pickers/DatePicker"
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider"
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs"
import dayjs, { Dayjs } from "dayjs"
import "dayjs/locale/ru"
import Card from "../Card"
import Button from "../Button"
import Input from "../Input"
import { SalesRecordFormData } from "../../libs/types"
import { createSalesRecord } from "../../libs/api/sales"

dayjs.locale("ru")

interface AddItemFormProps {
  onSuccess?: () => void
}

export default function AddItemForm({ onSuccess }: AddItemFormProps) {
  const [formData, setFormData] = useState<SalesRecordFormData>({
    type: "income",
    category: "",
    date: "",
    amount: 0,
  })

  const [selectedDate, setSelectedDate] = useState<Dayjs | null>(null)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setSuccess(false)
    setLoading(true)

    try {
      if (!formData.category.trim()) {
        setError("Категория обязательна для заполнения")
        setLoading(false)
        return
      }

      if (!selectedDate) {
        setError("Дата обязательна для заполнения")
        setLoading(false)
        return
      }

      if (formData.amount <= 0) {
        setError("Сумма должна быть положительным числом")
        setLoading(false)
        return
      }

      const apiData = {
        ...formData,
        date: selectedDate.format("YYYY-MM-DD"),
      }

      await createSalesRecord(apiData)

      setSuccess(true)

      setFormData({
        type: "income",
        category: "",
        date: "",
        amount: 0,
      })
      setSelectedDate(null)

      if (onSuccess) {
        onSuccess()
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка при добавлении записи"
      )
    } finally {
      setLoading(false)
    }
  }

  const handleChange =
    (field: keyof SalesRecordFormData) =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value =
        field === "amount" ? parseFloat(e.target.value) || 0 : e.target.value
      setFormData((prev) => ({ ...prev, [field]: value }))
    }

  return (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale="ru">
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
            label="Категория"
            value={formData.category}
            onChange={handleChange("category")}
            required
            placeholder="Введите категорию"
          />

          <DatePicker
            label="Дата"
            value={selectedDate}
            onChange={setSelectedDate}
            slotProps={{
              textField: {
                fullWidth: true,
                required: true,
                margin: "normal",
              },
            }}
          />

          <Input
            label="Сумма"
            type="number"
            value={formData.amount || ""}
            onChange={handleChange("amount")}
            required
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

          <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
            <Button type="submit" disabled={loading} sx={{ flex: 1 }}>
              {loading ? "Добавление..." : "Добавить запись"}
            </Button>
            <Link href="/" style={{ textDecoration: "none", flex: 1 }}>
              <Button variant="outlined" sx={{ width: "100%" }}>
                Отмена
              </Button>
            </Link>
          </Box>
        </Box>
      </Card>
    </LocalizationProvider>
  )
}
