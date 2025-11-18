"use client";

import { useState } from "react";
import { Container, Stack } from "@mui/material";
import Header from "./ui/Header";
import ImageUploadForm from "./ui/ImageUploadForm";
import ImageGallery from "./ui/ImageGallery";

export default function Home() {
  const [imageIds, setImageIds] = useState<string[]>([]);

  const handleUpload = (newImage: { id: string; status: string }) => {
    setImageIds((prev) => [...prev, newImage.id]);
  };

  const handleRemove = (id: string) => {
    setImageIds((prev) => prev.filter(imgId => imgId !== id));
  };

  return (
    <Container
      disableGutters
      sx={{
        width: "100%",
        maxWidth: 800,
        mx: "auto",
        px: { xs: 0, sm: 2 },
        py: 2,
        bgcolor: "background.default",
      }}
    >
      <Stack spacing={4}>
        <Header />
        <ImageUploadForm onUpload={handleUpload} />
        <ImageGallery imageIds={imageIds} onRemove={handleRemove} />
      </Stack>
    </Container>
  );
}