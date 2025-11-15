"use client";

import { useState } from "react";
import { useForm, Controller } from "react-hook-form";
import {
  Paper,
  Stack,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  IconButton,
  Link,
} from "@mui/material";
import { FieldRow } from "./FieldRow";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";

interface URLShortenerFormValues {
  originalUrl: string;
  customCode: string;
}

export default function URLShortenerForm() {
  const [result, setResult] = useState<{ shortUrl: string; originalUrl: string } | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<URLShortenerFormValues>({
    defaultValues: {
      originalUrl: "",
      customCode: "",
    },
  });

  const onSubmit = async (data: URLShortenerFormValues) => {
    setIsSubmitting(true);
    setError(null);
    setResult(null);

    const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

    try {
      const payload: { original: string; short?: string } = {
        original: data.originalUrl,
      };

      if (data.customCode.trim()) {
        payload.short = data.customCode.trim();
      }

      const response = await fetch(`${apiBase}/shorten`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `HTTP error: ${response.status}`);
      }

      const resData = await response.json();
      const shortUrl = resData.short;
      setResult({ shortUrl, originalUrl: data.originalUrl });
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Неизвестная ошибка";
      setError(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const getShortCode = (fullUrl: string): string => {
  try {
    const url = new URL(fullUrl);
    return url.pathname.split('/').pop() || '';
  } catch {
    return '';
  }
};

  return (
    <Paper>
      <Typography variant="h2" sx={{ textAlign: "center", color: "primary.dark", mb: 3 }}>
        Сократить URL
      </Typography>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Stack spacing={2}>
          <Controller
            name="originalUrl"
            control={form.control}
            rules={{ required: "URL обязателен", pattern: { value: /^https?:\/\//, message: "Введите корректный URL (с http:// или https://)" } }}
            render={({ field, fieldState }) => (
              <TextField
                {...field}
                label="Оригинальный URL"
                fullWidth
                error={!!fieldState.error}
                helperText={fieldState.error?.message}
              />
            )}
          />
          <Controller
            name="customCode"
            control={form.control}
            render={({ field }) => (
              <TextField
                {...field}
                label="Кастомный код (необязательно)"
                fullWidth
                placeholder="Например: my-link"
              />
            )}
          />
          <Button type="submit" variant="contained" disabled={isSubmitting}>
            {isSubmitting ? "Создание..." : "Создать короткую ссылку"}
          </Button>
        </Stack>
      </form>

      {error && (
        <Alert severity="error" sx={{ mt: 2 }}>
          {error}
        </Alert>
      )}

      {result && (
        <Box sx={{ mt: 3 }}>
          <Typography variant="h6" gutterBottom>
            Результат:
          </Typography>
          <FieldRow label="Оригинальный URL" value={result.originalUrl} />
          <FieldRow
            label="Короткая ссылка"
            value={
              <Box>
                <Link href={result.shortUrl} target="_blank" rel="noopener noreferrer" sx={{ display: "block", mb: 1 }}>
                  {result.shortUrl}
                </Link>
                <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                  <IconButton size="small" onClick={() => handleCopy(result.shortUrl)}>
                    <ContentCopyIcon fontSize="small" />
                  </IconButton>
                  <IconButton size="small" href={result.shortUrl} target="_blank" rel="noopener noreferrer">
                    <OpenInNewIcon fontSize="small" />
                  </IconButton>
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => handleCopy(getShortCode(result.shortUrl))}
                  >
                    Скопировать код
                  </Button>
                </Box>
              </Box>
            }
          />
        </Box>
      )}
    </Paper>
  );
}