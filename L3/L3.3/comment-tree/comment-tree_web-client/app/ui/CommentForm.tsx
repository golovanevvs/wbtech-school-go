"use client";

import { useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { Paper, Stack, TextField, Button, Alert } from "@mui/material";
import { useCommentContext } from "@/app/lib/CommentContext";

interface CommentFormValues {
  text: string;
}

export default function CommentForm({ parentId }: { parentId: number | null }) {
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { fetchComments } = useCommentContext();

  const form = useForm<CommentFormValues>({
    defaultValues: { text: "" },
  });

  const onSubmit = async (data: CommentFormValues) => {
    setIsSubmitting(true);
    setError(null);

    const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

    try {
      const payload: { text: string; parent_id?: number } = {
        text: data.text,
      };
      if (parentId !== null) {
        payload.parent_id = parentId;
      }

      const response = await fetch(`${apiBase}/comments`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `HTTP error: ${response.status}`);
      }

      form.reset();
      await fetchComments(); // Обновляем дерево
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Неизвестная ошибка";
      setError(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Paper sx={{ p: 2 }}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Stack spacing={2}>
          <Controller
            name="text"
            control={form.control}
            rules={{ required: "Текст комментария обязателен" }}
            render={({ field, fieldState }) => (
              <TextField
                {...field}
                label={parentId ? "Ответить" : "Добавить комментарий"}
                multiline
                rows={3}
                fullWidth
                error={!!fieldState.error}
                helperText={fieldState.error?.message}
              />
            )}
          />
          <Button type="submit" variant="contained" disabled={isSubmitting}>
            {isSubmitting ? "Отправка..." : parentId ? "Ответить" : "Добавить"}
          </Button>
        </Stack>
      </form>

      {error && (
        <Alert severity="error" sx={{ mt: 2 }}>
          {error}
        </Alert>
      )}
    </Paper>
  );
}