"use client"

import { useState, useEffect, useRef } from "react"
import {
  Box,
  Card,
  CardMedia,
  CardContent,
  CardActions,
  Typography,
  IconButton,
} from "@mui/material"
import DeleteIcon from "@mui/icons-material/Delete"
import { getImageStatus, deleteImage } from "../lib/api"
import { Image } from "../lib/types"

interface Props {
  id: number // ✅ id как number
  onRemove: (id: number) => void // ✅ onRemove как (id: number) => void
}

export default function ImageCard({ id, onRemove }: Props) {
  const [image, setImage] = useState<Image | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const isMountedRef = useRef(true)
  const timerRef = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    const fetchStatus = async () => {
      console.log("fetchStatus called for id:", id)
      if (!isMountedRef.current) return

      try {
        const data = await getImageStatus(id.toString()) // ✅ преобразуем к string для API
        console.log("getImageStatus response received:", data)

        if (!isMountedRef.current) {
          console.log("Component unmounted, exiting")
          return
        }

        // Проверим тип вручную
        console.log("Checking data types...")
        if (typeof data.id !== "number" || typeof data.status !== "string") {
          console.error("Invalid data format:", data)
          setError("Invalid data format")
          return
        }
        console.log("Data types are valid")

        console.log("Before setImage:", { id: data.id, status: data.status })

        setImage(data)
        console.log("After setImage")

        if (data.status === "completed" || data.status === "failed") {
          console.log("Status is final, not polling")
          return
        }

        timerRef.current = setTimeout(fetchStatus, 2000)
      } catch (err: unknown) {
        console.log("Error in fetchStatus:", err)
        if (!isMountedRef.current) return

        let message = "Не удалось получить статус изображения"
        if (err instanceof Error) {
          message = err.message
        }
        setError(message)
      }
    }

    fetchStatus()

    return () => {
      isMountedRef.current = false
      if (timerRef.current) {
        clearTimeout(timerRef.current)
      }
    }
  }, [id])

  console.log("ImageCard render, image:", image, "error:", error)

  const handleDelete = async () => {
    setLoading(true)
    try {
      await deleteImage(id.toString()) // ✅ преобразуем к string для API
      onRemove(id)
    } catch (err: unknown) {
      let message = "Ошибка удаления"
      if (err instanceof Error) {
        message = err.message
      }
      alert(message)
    } finally {
      setLoading(false)
    }
  }

  if (error) return <Box color="error.main">❌ Ошибка: {error}</Box>

  if (!image) {
    console.log("Image is null, showing loading")
    return <Box>⏳ Загрузка...</Box>
  }

  console.log("Image is not null, showing card")

  return (
    <Card sx={{ maxWidth: 345, margin: "10px" }}>
      <CardContent>
        <Typography gutterBottom variant="h6" component="div">
          ID изображения: {image.id}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Статус: {image.status}
        </Typography>
      </CardContent>

      {image.status === "completed" && image.processed_url && (
        <CardMedia
          component="img"
          height="140"
          image={image.processed_url}
          alt="Обработанное изображение"
        />
      )}

      {image.status === "failed" && (
        <CardContent>
          <Typography variant="body2" color="error.main">
            ❌ Ошибка обработки
          </Typography>
        </CardContent>
      )}

      {image.status === "processing" && (
        <CardContent>
          <Typography variant="body2" color="text.secondary">
            🔄 Обработка...
          </Typography>
        </CardContent>
      )}

      <CardActions>
        <IconButton
          size="small"
          color="error"
          onClick={handleDelete}
          disabled={loading}
          sx={{ ml: "auto" }}
        >
          <DeleteIcon />
        </IconButton>
      </CardActions>
    </Card>
  )
}