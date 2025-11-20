"use client"

import { Box } from "@mui/material"
import ImageCard from "./ImageCard"

export default function ImageGallery({
  imageIds,
  onRemove,
}: {
  imageIds: number[]
  onRemove: (id: number) => void
}) {
  return (
    <Box
      sx={{
        width: "100%",
      }}
    >
      <Box
        sx={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(310px, 1fr))",
          gap: 2,
          width: "fit-content",
          maxWidth: "100%",
          mx: "auto",
        }}
      >
        {imageIds.map((id) => (
          <ImageCard key={id} id={id} onRemove={onRemove} />
        ))}
      </Box>
    </Box>
  )
}
