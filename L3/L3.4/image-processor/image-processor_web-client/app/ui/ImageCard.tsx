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
import DownloadIcon from "@mui/icons-material/Download"
import { getImageStatus, deleteImage } from "../lib/api"
import { Image } from "../lib/types"

interface Props {
  id: number
  onRemove: (id: number) => void
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
        const data = await getImageStatus(id.toString())
        console.log("getImageStatus response received:", data)

        if (!isMountedRef.current) {
          console.log("Component unmounted, exiting")
          return
        }

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
      await deleteImage(id.toString())
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

  const handleDownload = () => {
    if (image?.processed_url) {
      window.open(image.processed_url, "_blank")
    }
  }

  if (error) return <Box color="error.main">❌ Ошибка: {error}</Box>

  if (!image) {
    console.log("Image is null, showing loading")
    return <Box>⏳ Загрузка...</Box>
  }

  console.log("Image is not null, showing card")

  const originalFilename = image.original_path
    ? image.original_path.split("\\").pop()?.split("/").pop() || "unknown"
    : "unknown"

  const operationsList = []
  if (image.operations) {
    if (image.operations.resize) operationsList.push("изменение размера")
    if (image.operations.watermark)
      operationsList.push("наложение водяного знака")
    if (image.operations.thumbnail) operationsList.push("создание миниатюры")
  }
  const operationsText =
    operationsList.length > 0 ? operationsList.join(", ") : "без изменений"

  return (
    <Card
      sx={{
        width: 320,
        margin: "10px",
        display: "flex",
        flexDirection: "column",
        height: "100%",
      }}
    >
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography gutterBottom variant="h6" component="div">
          ID изображения: {image.id}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Статус: {image.status}
        </Typography>
        <Typography
          variant="body2"
          color="text.secondary"
          sx={{
            wordBreak: "break-word",
            overflowWrap: "break-word",
          }}
        >
          Файл: {originalFilename}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Обработка: {operationsText}
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

      {(image.status === "uploading" || image.status === "processing") && (
        <CardContent
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: 140,
          }}
        >
          <Typography variant="body2" color="text.secondary">
            📁 Загрузка...
          </Typography>
        </CardContent>
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
        {image.status === "completed" && image.processed_url && (
          <IconButton
            size="small"
            color="primary"
            onClick={handleDownload}
            aria-label="открыть в новом окне"
          >
            <DownloadIcon />
          </IconButton>
        )}
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
