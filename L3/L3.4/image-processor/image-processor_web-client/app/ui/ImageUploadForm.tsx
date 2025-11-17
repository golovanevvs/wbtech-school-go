"use client";

import { useState, useRef, ChangeEvent } from "react";
import {
  Paper,
  Button,
  Typography,
} from "@mui/material";
import { uploadImage } from "../lib/api";

interface Props {
  onUpload: (image: { id: string; status: string }) => void;
}

export default function ImageUploadForm({ onUpload }: Props) {
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];
      setFile(selectedFile);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    setLoading(true);
    try {
      const { id } = await uploadImage(file);
      onUpload({ id, status: "uploading" });
      setFile(null); // Reset input
      if (fileInputRef.current) fileInputRef.current.value = ""; // Clear input
    } catch (err: unknown) {
      let message = "Upload failed";
      if (err instanceof Error) {
        message = err.message;
      }
      alert(message);
    }
    setLoading(false);
  };

  return (
    <Paper
      component="form"
      onSubmit={handleSubmit}
      sx={{ p: 2, display: "flex", gap: 2, alignItems: "center", flexDirection: "column" }}
    >
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/*"
        style={{ display: "none" }}
      />
      <Button
        variant="outlined"
        onClick={handleButtonClick}
        disabled={loading}
      >
        Выберите файл
      </Button>
      {file && (
        <Typography variant="body2">
          Выбрано: {file.name}
        </Typography>
      )}
      <Button
        type="submit"
        variant="contained"
        disabled={loading || !file}
      >
        {loading ? 'Загрузка...' : 'Загрузить'}
      </Button>
    </Paper>
  );
}