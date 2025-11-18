"use client";

import { useState, useEffect, useRef } from "react";
import { Box, Card, CardMedia, CardContent, CardActions, Typography, IconButton } from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import { getImageStatus, deleteImage } from "../lib/api";
import { Image } from "../lib/types";

interface Props {
  id: string;
  onRemove: (id: string) => void;
}

export default function ImageCard({ id, onRemove }: Props) {
  const [image, setImage] = useState<Image | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const isMountedRef = useRef(true);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const fetchStatus = async () => {
      if (!isMountedRef.current) return;

      try {
        const data = await getImageStatus(id);
        if (!isMountedRef.current) return;

        setImage(data);

        if (data.status === "completed" || data.status === "failed") {
          return;
        }

        timerRef.current = setTimeout(fetchStatus, 2000);
      } catch (err: unknown) {
        if (!isMountedRef.current) return;

        let message = "Не удалось получить статус изображения";
        if (err instanceof Error) {
          message = err.message;
        }
        setError(message);
      }
    };

    fetchStatus();

    return () => {
      isMountedRef.current = false;
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, [id]);

  const handleDelete = async () => {
    setLoading(true);
    try {
      await deleteImage(id);
      onRemove(id);
    } catch (err: unknown) {
      let message = "Ошибка удаления";
      if (err instanceof Error) {
        message = err.message;
      }
      alert(message);
    } finally {
      setLoading(false);
    }
  };

  if (error) return <Box color="error.main">❌ Ошибка: {error}</Box>;

  if (!image) return <Box>⏳ Загрузка...</Box>;

  return (
    <Card sx={{ maxWidth: 345, margin: "10px" }}>
      <CardContent>
        <Typography gutterBottom variant="h6" component="div">
          ID изображения: {image.id}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Статус: {image.status}
        </Typography>
      </CardContent>

      {image.status === "completed" && image.processedUrl && (
        <CardMedia
          component="img"
          height="140"
          image={image.processedUrl}
          alt="Обработанное изображение"
        />
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
  );
}