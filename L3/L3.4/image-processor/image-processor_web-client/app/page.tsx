"use client"

import { useState, useEffect } from "react"
import { Box, Stack } from "@mui/material"
import Header from "./ui/Header"
import ImageUploadForm from "./ui/ImageUploadForm"
import ImageGallery from "./ui/ImageGallery"
import { getAllImages } from "./lib/api"

export default function Home() {
  const [imageIds, setImageIds] = useState<number[]>([])

  const handleUpload = (newImage: { id: number; status: string }) => {
    setImageIds((prev) => [newImage.id, ...prev])
  }

  const handleRemove = (id: number) => {
    setImageIds((prev) => prev.filter((imgId) => imgId !== id))
  }

  useEffect(() => {
    const fetchAllImages = async () => {
      try {
        const images = await getAllImages()
        const ids = images.map((img) => img.id)
        setImageIds(ids)
      } catch (err) {
        console.error("Failed to load images:", err)
      }
    }

    fetchAllImages()
  }, [])

  return (
    <Box
      sx={{
        width: "100%",
        minHeight: "100vh",
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
        maxWidth: "100vw", 
        mx: "auto", 
      }}
    >
      <Stack spacing={4}>
        <Header />
        <ImageUploadForm onUpload={handleUpload} />
        <ImageGallery imageIds={imageIds} onRemove={handleRemove} />
      </Stack>
    </Box>
  )
}
