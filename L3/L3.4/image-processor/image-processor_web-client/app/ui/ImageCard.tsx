"use client";

import { useState, useEffect } from "react";
import {
  Box,
  Card,
  CardMedia,
  CardContent,
  CardActions,
  Button,
  Typography,
} from "@mui/material";
import { getImageStatus, deleteImage } from "../lib/api";
import { Image } from "../lib/types";

interface Props {
  id: string;
}

export default function ImageCard({ id }: Props) {
  const [image, setImage] = useState<Image | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const data = await getImageStatus(id);
        setImage(data);
        if (data.status === "completed" || data.status === "failed") return;
        // Poll every 2 seconds
        setTimeout(fetchStatus, 2000);
      } catch (err: unknown) {
        let message = "Failed to fetch image status";
        if (err instanceof Error) {
          message = err.message;
        }
        setError(message);
      }
    };

    fetchStatus();
  }, [id]);

  if (error) return <Box color="error.main">❌ Error: {error}</Box>;

  if (!image) return <Box>⏳ Loading...</Box>;

  const handleDelete = async () => {
    try {
      await deleteImage(id);
      // Optionally, notify parent to remove from list
    } catch (err: unknown) {
      let message = "Delete failed";
      if (err instanceof Error) {
        message = err.message;
      }
      alert(message);
    }
  };

  return (
    <Card sx={{ maxWidth: 345, margin: "10px" }}>
      <CardContent>
        <Typography gutterBottom variant="h6" component="div">
          Image ID: {image.id}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Status: {image.status}
        </Typography>
      </CardContent>
      {image.status === "completed" && image.processedUrl && (
        <CardMedia
          component="img"
          height="140"
          image={image.processedUrl}
          alt="Processed"
        />
      )}
      {image.status === "processing" && (
        <CardContent>
          <Typography variant="body2" color="text.secondary">
            🔄 Processing...
          </Typography>
        </CardContent>
      )}
      <CardActions>
        <Button size="small" color="error" onClick={handleDelete}>
          Delete
        </Button>
      </CardActions>
    </Card>
  );
}