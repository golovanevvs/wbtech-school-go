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
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
        gap: 2,
        justifyContent: 'center',
        width: '100%',
      }}
    >
      {imageIds.map((id) => (
        <ImageCard 
          key={id} 
          id={id}
          onRemove={onRemove}
        />
      ))}
    </Box>
  )
}