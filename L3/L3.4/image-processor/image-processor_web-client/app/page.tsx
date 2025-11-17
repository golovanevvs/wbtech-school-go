"use client";

import { useState } from "react";
import { Container, Stack } from "@mui/material";
import Header from "./ui/Header";
import ImageUploadForm from "./ui/ImageUploadForm";
import ImageGallery from "./ui/ImageGallery";
import { Image } from "./lib/types";

export default function Home() {
  const [images, setImages] = useState<Image[]>([]);

  const handleUpload = (newImage: { id: string; status: string }) => {
    setImages((prev) => [...prev, { id: newImage.id, status: newImage.status } as Image]);
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
        <ImageGallery initialImages={images} />
      </Stack>
    </Container>
  );
}