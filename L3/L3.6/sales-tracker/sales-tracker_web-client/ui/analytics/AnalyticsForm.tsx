"use client"

import { useState } from "react"
import Link from "next/link"
import { Box, Typography, Alert, Paper } from "@mui/material"
import { DatePicker } from "@mui/x-date-pickers/DatePicker"
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider"
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs"
import dayjs, { Dayjs } from "dayjs"
import "dayjs/locale/ru"
import Card from "../Card"
import Button from "../Button"
import { AnalyticsData, AnalyticsRequest } from "../../libs/types"
import { getAnalytics } from "../../libs/api/sales"

dayjs.locale("ru")

export default function AnalyticsForm() {
  const [fromDate, setFromDate] = useState<Dayjs | null>(
    dayjs().subtract(30, "day")
  )
  const [toDate, setToDate] = useState<Dayjs | null>(dayjs())
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const handleGetAnalytics = async () => {
    if (!fromDate || !toDate) {
      setError("Пожалуйста, выберите даты начала и окончания")
      return
    }

    if (fromDate.isAfter(toDate)) {
      setError("Дата начала не может быть позже даты окончания")
      return
    }

    setError("")
    setLoading(true)

    try {
      const request: AnalyticsRequest = {
        from: fromDate.startOf("day").toISOString(),
        to: toDate.endOf("day").toISOString(),
      }

      const data = await getAnalytics(request)
      setAnalytics(data)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка при получении аналитики"
      )
      setAnalytics(null)
    } finally {
      setLoading(false)
    }
  }

  const formatNumber = (num: number) => {
    return num.toLocaleString("ru-RU", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })
  }

  const formatCurrency = (num: number) => {
    return `${formatNumber(num)} ₽`
  }

  return (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale="ru">
      <Box sx={{ width: "100%", maxWidth: 800, mx: "auto" }}>
        <Typography variant="h4" component="h1" gutterBottom align="center">
          Аналитика продаж
        </Typography>

        <Card>
          <Typography variant="h6" component="h2" gutterBottom>
            Выберите период анализа
          </Typography>

          <Box sx={{ display: "flex", gap: 2, mb: 3, flexWrap: "wrap" }}>
            <DatePicker
              label="Дата начала"
              value={fromDate}
              onChange={setFromDate}
              slotProps={{
                textField: {
                  fullWidth: true,
                  required: true,
                },
              }}
            />

            <DatePicker
              label="Дата окончания"
              value={toDate}
              onChange={setToDate}
              slotProps={{
                textField: {
                  fullWidth: true,
                  required: true,
                },
              }}
            />
          </Box>

          <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
            <Button
              onClick={handleGetAnalytics}
              disabled={loading || !fromDate || !toDate}
              sx={{ mb: 3 }}
            >
              {loading ? "Загрузка..." : "Получить аналитику"}
            </Button>
            <Link href="/" style={{ textDecoration: "none", flex: 1 }}>
              <Button variant="outlined" sx={{ width: "100%" }}>
                Отмена
              </Button>
            </Link>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {analytics && (
            <Box sx={{ mt: 3 }}>
              <Typography variant="h6" component="h3" gutterBottom>
                Результаты анализа за период:
                <br />
                {fromDate?.format("DD.MM.YYYY")} -{" "}
                {toDate?.format("DD.MM.YYYY")}
              </Typography>

              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))",
                  gap: 2,
                }}
              >
                <Paper sx={{ p: 2, textAlign: "center" }}>
                  <Typography variant="h4" color="primary" gutterBottom>
                    {formatCurrency(analytics.sum)}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    Общая сумма
                  </Typography>
                </Paper>

                <Paper sx={{ p: 2, textAlign: "center" }}>
                  <Typography variant="h4" color="secondary" gutterBottom>
                    {formatCurrency(analytics.avg)}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    Средний чек
                  </Typography>
                </Paper>

                <Paper sx={{ p: 2, textAlign: "center" }}>
                  <Typography variant="h4" color="info.main" gutterBottom>
                    {analytics.count}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    Количество записей
                  </Typography>
                </Paper>

                <Paper sx={{ p: 2, textAlign: "center" }}>
                  <Typography variant="h4" color="success.main" gutterBottom>
                    {formatCurrency(analytics.median)}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    Медиана
                  </Typography>
                </Paper>

                <Paper sx={{ p: 2, textAlign: "center" }}>
                  <Typography variant="h4" color="warning.main" gutterBottom>
                    {formatCurrency(analytics.percentile90)}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    90-й перцентиль
                  </Typography>
                </Paper>
              </Box>
            </Box>
          )}
        </Card>
      </Box>
    </LocalizationProvider>
  )
}
