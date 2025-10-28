import IconButton from "@mui/material/IconButton"
import InputAdornment from "@mui/material/InputAdornment"
import Paper from "@mui/material/Paper"
import Stack from "@mui/material/Stack"
import TextField from "@mui/material/TextField"
import { useState } from "react"
import { FieldRow } from "./FieldRow"
import SearchIcon from "@mui/icons-material/Search"

export default function GetNoticeStatus() {
  const [statusID, setStatusID] = useState("")
  const [status, setStatus] = useState("Ожидание запроса")

  const statusColor = (s: string) => {
    if (s.includes("Ошибка")) return "error.main"
    if (s.includes("Успешно")) return "success.main"
    return "primary.main"
  }

  const onGetStatus = async () => {
    if (!statusID) {
      setStatus("Ошибка: укажите ID уведомления")
      return
    }

    setStatus("Получение статуса...")

    try {
      const response = await fetch(
        `https://insecurely-fond-shiner.cloudpub.ru/notify/${statusID}`
      )
      const data = await response.json()

      if (!response.ok)
        throw new Error(data.error || `HTTP error: ${response.status}`)

      setStatus(`Статус: ${data.status ?? "неизвестен"}`)
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Неизвестная ошибка"
      setStatus(`Ошибка: ${msg}`)
    }
  }
  return (
    <Paper sx={{ p: 3 }}>
      <Stack spacing={2}>
        <TextField
          label="ID уведомления"
          value={statusID}
          onChange={(e) => setStatusID(e.target.value)}
          fullWidth
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton onClick={onGetStatus}>
                  <SearchIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
        <FieldRow
          label="Статус"
          value={status}
          statusColor={statusColor(status)}
        />
      </Stack>
    </Paper>
  )
}
