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
  Chip,
} from "@mui/material";
import { AnalyticsData } from "@/app/lib/types";

interface AnalyticsFormValues {
  shortCode: string;
}

export default function AnalyticsSection() {
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const form = useForm<AnalyticsFormValues>({
    defaultValues: {
      shortCode: "",
    },
  });

  const onSubmit = async (data: AnalyticsFormValues) => {
    setLoading(true);
    setError(null);
    setAnalytics(null);

    const apiBase = process.env.NEXT_PUBLIC_API_URL;
    const shortCode = data.shortCode.trim();

    if (!shortCode) {
      setError("Код короткой ссылки не может быть пустым");
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`${apiBase}/analytics/${shortCode}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `HTTP error: ${response.status}`);
      }

      const resData: AnalyticsData = await response.json();
      setAnalytics(resData);
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Неизвестная ошибка";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper>
      <Typography variant="h2" sx={{ textAlign: "center", color: "primary.dark", mb: 3 }}>
        Аналитика по ссылке
      </Typography>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Stack spacing={2}>
          <Controller
            name="shortCode"
            control={form.control}
            rules={{ required: "Код обязателен" }}
            render={({ field, fieldState }) => (
              <TextField
                {...field}
                label="Код короткой ссылки"
                fullWidth
                placeholder="Например: abc123"
                error={!!fieldState.error}
                helperText={fieldState.error?.message}
              />
            )}
          />
          <Button type="submit" variant="outlined" disabled={loading}>
            {loading ? "Загрузка..." : "Получить аналитику"}
          </Button>
        </Stack>
      </form>

      {error && (
        <Alert severity="error" sx={{ mt: 2 }}>
          {error}
        </Alert>
      )}

      {analytics && (
        <Box sx={{ mt: 3 }}>
          <Typography variant="h6" gutterBottom>
            Всего переходов: <Chip label={analytics.total_clicks} color="primary" />
          </Typography>
          {analytics.clicks && Array.isArray(analytics.clicks) && analytics.clicks.length > 0 ? (
            <Stack spacing={2}>
              {analytics.clicks.map((click) => (
                <Paper key={click.id} sx={{ p: 2 }}>
                  <Box sx={{ mb: 1 }}>
                    <Typography variant="body2" color="text.secondary">Дата и время</Typography>
                    <Typography variant="body1">{new Date(click.created_at).toLocaleString()}</Typography>
                  </Box>
                  <Box sx={{ mb: 1 }}>
                    <Typography variant="body2" color="text.secondary">User Agent</Typography>
                    <Typography variant="body1" sx={{
                      wordBreak: "break-word",
                      overflowWrap: "break-word",
                    }}>
                      {click.user_agent}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 1 }}>
                    <Typography variant="body2" color="text.secondary">IP</Typography>
                    <Typography variant="body1">{click.ip || "Неизвестен"}</Typography>
                  </Box>
                </Paper>
              ))}
            </Stack>
          ) : (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              Нет данных о переходах.
            </Typography>
          )}
        </Box>
      )}
    </Paper>
  );
}