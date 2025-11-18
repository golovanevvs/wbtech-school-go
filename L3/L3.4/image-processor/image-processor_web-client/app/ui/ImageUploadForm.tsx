"use client"

import { useState, useRef, ChangeEvent } from "react"
import {
  Paper,
  Button,
  Stack,
  Typography,
  FormControlLabel,
  Checkbox,
} from "@mui/material"
import { uploadImage } from "../lib/api"

interface Props {
  onUpload: (image: { id: string; status: string }) => void
}

export default function ImageUploadForm({ onUpload }: Props) {
  const [file, setFile] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)

  const [resize, setResize] = useState(true)
  const [thumbnail, setThumbnail] = useState(false)
  const [watermark, setWatermark] = useState(false)

  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleButtonClick = () => {
    fileInputRef.current?.click()
  }

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0]
      setFile(selectedFile)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!file) return

    setLoading(true)
    try {
      const { id } = await uploadImage(file, { resize, thumbnail, watermark })
      onUpload({ id, status: "uploading" })
      setFile(null)
      if (fileInputRef.current) fileInputRef.current.value = ""
    } catch (err: unknown) {
      let message = "Upload failed"
      if (err instanceof Error) {
        message = err.message
      }
      alert(message)
    }
    setLoading(false)
  }

  return (
    <Paper
      component="form"
      onSubmit={handleSubmit}
      sx={{ p: 2, display: "flex", gap: 2, flexDirection: "column" }}
    >
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/*"
        style={{ display: "none" }}
      />
      <Button variant="outlined" onClick={handleButtonClick} disabled={loading}>
        Выберите файл
      </Button>
      {file && <Typography variant="body2">Выбрано: {file.name}</Typography>}

      <Stack spacing={1}>
        <FormControlLabel
          control={
            <Checkbox
              checked={resize}
              onChange={(e) => setResize(e.target.checked)}
              disabled={!file || loading}
            />
          }
          label="Изменить размер"
        />
        <FormControlLabel
          control={
            <Checkbox
              checked={thumbnail}
              onChange={(e) => setThumbnail(e.target.checked)}
              disabled={!file || loading}
            />
          }
          label="Создать миниатюру"
        />
        <FormControlLabel
          control={
            <Checkbox
              checked={watermark}
              onChange={(e) => setWatermark(e.target.checked)}
              disabled={!file || loading}
            />
          }
          label="Наложить водяной знак"
        />
      </Stack>

      <Button type="submit" variant="contained" disabled={loading || !file}>
        {loading ? "Загрузка..." : "Загрузить"}
      </Button>
    </Paper>
  )
}
