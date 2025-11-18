"use client";

import { Stack } from "@mui/material";
import ImageCard from "./ImageCard";

export default function ImageGallery({
  imageIds,
  onRemove,
}: {
  imageIds: string[];
  onRemove: (id: string) => void;
}) {
  return (
    <Stack spacing={2} direction="row" flexWrap="wrap" justifyContent="center">
      {imageIds.map((id) => (
        <ImageCard key={id} id={id} onRemove={onRemove} />
      ))}
    </Stack>
  );
}