"use client"

import { useForm, Controller } from "react-hook-form"
import TextField from "@mui/material/TextField"
import { Paper } from "@mui/material"
import InputAdornment from "@mui/material/InputAdornment"
import IconButton from "@mui/material/IconButton"
import SearchIcon from "@mui/icons-material/Search"
import { OrderResponse } from "@/app/lib/types"
import { FieldRow } from "@/app/ui/FieldRow"
import { useState } from "react"

interface HookFormProps {
  onOrderDataReceived: (response: OrderResponse) => void
}

interface FormValues {
  orderUID: string
}

export default function HookForm({ onOrderDataReceived }: HookFormProps) {
  const { control, handleSubmit, reset } = useForm<FormValues>({
    defaultValues: {
      orderUID: "",
    },
  })

  const [requestStatus, setRequestStatus] = useState<string>("Ожидание запроса")

  const getStatusColor = () => {
    if (requestStatus === "Ожидание запроса") return "primary.main"
    if (requestStatus.includes("Ошибка")) return "error.main"
    if (requestStatus === "Данные успешно получены") return "success.main"
    return "text.secondary"
  }

  const onSubmit = async (formData: FormValues) => {
    setRequestStatus("Ожидание запроса")

    await new Promise((resolve) => setTimeout(resolve, 50))

    setRequestStatus("Запрос отправляется...")

    try {
      const response = await fetch(`/order?order_uid=${formData.orderUID}`, {
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      })

      if (!response.ok) {
        let errorMessage = `HTTP error! status: ${response.status}`

        try {
          const data = await response.json()
          if (data?.error) {
            errorMessage = data.error
          }
        } catch {
        }

        setRequestStatus(`Ошибка: ${errorMessage}`)
        onOrderDataReceived({
          success: false,
          error: errorMessage,
        })
        return
      }

      const data = await response.json()

      setRequestStatus("Данные успешно получены")
      onOrderDataReceived({
        success: true,
        data,
      })
    } catch (error) {
      let errorMessage: string

      if (error instanceof TypeError && error.message === "Failed to fetch") {
        errorMessage = "Нет связи с сервером"
      } else {
        errorMessage = error instanceof Error ? error.message : "Неизвестная ошибка"
      }

      setRequestStatus(`Ошибка: ${errorMessage}`)
      console.error("Ошибка при получении данных:", error)
      onOrderDataReceived({
        success: false,
        error: errorMessage,
      })
    }
  }

  return (
    <Paper sx={{ gap: "20px" }}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Controller
          name="orderUID"
          control={control}
          render={({ field }) => (
            <TextField
              {...field}
              fullWidth
              label="Введите UID заказа"
              variant="outlined"
              sx={{
                "& .MuiOutlinedInput-root": {
                  "& fieldset": {
                    borderColor: "primary.light",
                    borderWidth: 1,
                  },
                  "&:hover fieldset": {
                    borderColor: "primary.main",
                  },
                  "&.Mui-focused fieldset": {
                    borderWidth: 2,
                    borderColor: "primary.main",
                  },
                  backgroundColor: "background.paper",
                  color: "primary.dark",
                },
                "& .MuiInputLabel-root": {
                  color: "primary.light",
                },
                "& .MuiInputLabel-root.Mui-focused": {
                  color: "primary.main",
                },
              }}
              slotProps={{
                input: {
                  autoComplete: "off",
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        edge="end"
                        type="submit"
                        sx={{
                          color: "primary.light",
                          "&:hover": { color: "primary.main" },
                          "&:active": { color: "primary.dark" },
                        }}
                      >
                        <SearchIcon />
                      </IconButton>
                    </InputAdornment>
                  ),
                },
              }}
            />
          )}
        />
      </form>

      <FieldRow
        label="Статус запроса"
        value={requestStatus}
        statusColor={getStatusColor()}
      />
    </Paper>
  )
}
