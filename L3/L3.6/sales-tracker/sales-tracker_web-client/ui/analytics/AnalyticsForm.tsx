"use client"

import { useState } from "react"
import Link from "next/link"
import {
  Box,
  Typography,
  Alert,
  Paper,
  IconButton,
  Tooltip,
} from "@mui/material"
import { DatePicker } from "@mui/x-date-pickers/DatePicker"
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider"
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs"
import dayjs, { Dayjs } from "dayjs"
import "dayjs/locale/ru"
import DownloadIcon from "@mui/icons-material/Download"
import { LineChart } from "@mui/x-charts/LineChart"
import Card from "../Card"
import Button from "../Button"
import { AnalyticsData, AnalyticsRequest, SalesRecord } from "../../libs/types"
import {
  getAnalytics,
  getSalesRecords,
  downloadCSV,
} from "../../libs/api/sales"

dayjs.locale("ru")

export default function AnalyticsForm() {
  const [fromDate, setFromDate] = useState<Dayjs | null>(
    dayjs().subtract(30, "day")
  )
  const [toDate, setToDate] = useState<Dayjs | null>(dayjs())
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null)
  const [chartData, setChartData] = useState<{
    dates: string[]
    amounts: number[]
  }>({ dates: [], amounts: [] })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const aggregateDataByDay = (
    records: SalesRecord[]
  ): { dates: string[]; amounts: number[] } => {
    const dailyData: { [key: string]: number } = {}

    records.forEach((record) => {
      const date = dayjs(record.date).format("DD.MM")
      if (!dailyData[date]) {
        dailyData[date] = 0
      }
      if (record.type === "income") {
        dailyData[date] += record.amount
      } else {
        dailyData[date] -= record.amount
      }
    })

    const dates = Object.keys(dailyData).sort((a, b) => {
      const [dayA, monthA] = a.split(".").map(Number)
      const [dayB, monthB] = b.split(".").map(Number)
      return (
        new Date(2024, monthA - 1, dayA).getTime() -
        new Date(2024, monthB - 1, dayB).getTime()
      )
    })
    const amounts = dates.map((date) => dailyData[date])

    return { dates, amounts }
  }

  const handleDownloadCSV = async () => {
    try {
      if (!fromDate || !toDate) {
        setError("Пожалуйста, выберите даты для скачивания CSV")
        return
      }

      const blob = await downloadCSV(
        fromDate.startOf("day").toISOString(),
        toDate.endOf("day").toISOString()
      )

      const url = window.URL.createObjectURL(blob)
      const a = document.createElement("a")
      a.style.display = "none"
      a.href = url
      a.download = `sales-data-${fromDate.format(
        "YYYY-MM-DD"
      )}-to-${toDate.format("YYYY-MM-DD")}.csv`

      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)

      setError("")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка при скачивании CSV")
    }
  }

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

      const [analyticsData, salesRecords] = await Promise.all([
        getAnalytics(request),
        getSalesRecords(),
      ])

      setAnalytics(analyticsData)

      const filteredRecords = salesRecords.filter((record) => {
        const recordDate = dayjs(record.date)
        return (
          recordDate.isAfter(fromDate.startOf("day")) &&
          recordDate.isBefore(toDate.endOf("day"))
        )
      })

      const aggregatedData = aggregateDataByDay(filteredRecords)
      setChartData(aggregatedData)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка при получении аналитики"
      )
      setAnalytics(null)
      setChartData({ dates: [], amounts: [] })
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

          <Box sx={{ display: "flex", gap: 2, mt: 2, flexWrap: "wrap" }}>
            <Button
              onClick={handleGetAnalytics}
              disabled={loading || !fromDate || !toDate}
              sx={{ mb: 3 }}
            >
              {loading ? "Загрузка..." : "Получить аналитику"}
            </Button>

            <Link
              href="/"
              style={{ textDecoration: "none", flex: 1, minWidth: "120px" }}
            >
              <Button variant="outlined" sx={{ width: "100%", mb: 3 }}>
                Отмена
              </Button>
            </Link>
          </Box>
        </Card>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {analytics && (
          <Box
            sx={{
              mt: 1,
              mx: "auto",
            }}
          >
            <Typography variant="h6" component="h3" gutterBottom>
              Результаты анализа за период:
              <br />
              {fromDate?.format("DD.MM.YYYY")} - {toDate?.format("DD.MM.YYYY")}
            </Typography>

            <Tooltip title="Скачать CSV файл">
              <span>
                <IconButton
                  onClick={handleDownloadCSV}
                  disabled={!fromDate || !toDate}
                  color="primary"
                  sx={{ mb: 3 }}
                >
                  <DownloadIcon />
                </IconButton>
              </span>
            </Tooltip>

            <Box
              sx={{
                display: "grid",
                maxWidth: 300,
                gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))",
                gap: 1,
                mb: 1,
                mx: "auto",
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

            {/* График динамики продаж */}
            {chartData.dates.length > 0 && (
              <Box sx={{ mt: 2 }}>
                <Typography variant="h6" component="h3" gutterBottom>
                  Динамика продаж по дням
                </Typography>
                <Paper sx={{ p: 2 }}>
                  <LineChart
                    xAxis={[
                      {
                        scaleType: "point",
                        data: chartData.dates,
                        label: "Дата",
                      },
                    ]}
                    series={[
                      {
                        data: chartData.amounts,
                        label: "Сумма (₽)",
                        color: "#1976d2",
                      },
                    ]}
                    width={300}
                    height={150}
                    margin={{ left: 0, right: 5, top: 0, bottom: 0 }}
                    sx={{
                      "& .MuiLineChart-root": {
                        overflow: "visible",
                      },
                    }}
                  />
                </Paper>
              </Box>
            )}
          </Box>
        )}
      </Box>
    </LocalizationProvider>
  )
}
