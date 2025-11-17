"use client"

import { useState } from "react"
import { Stack } from "@mui/material"
import ImageCard from "./ImageCard"
import { Image } from "../lib/types"

interface Props {
  initialImages: Image[]
}

export default function ImageGallery({ initialImages }: Props) {
  const [images, setImages] = useState<Image[]>(initialImages)

  const handleImageAdded = (newImage: { id: string; status: string }) => {
    setImages((prev) => [
      ...prev,
      { id: newImage.id, status: newImage.status } as Image,
    ])
  }

  return (
    <Stack spacing={2} direction="row" flexWrap="wrap" justifyContent="center">
      {images.map((img) => (
        <ImageCard key={img.id} id={img.id} />
      ))}
    </Stack>
  )
}
