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
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Collapse,
  IconButton,
} from "@mui/material";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import { AnalyticsData } from "@/app/lib/types";

interface AnalyticsFormValues {
  shortCode: string;
}

function Row(props: { row: AnalyticsData['clicks'][0] }) {
  const { row } = props;
  const [open, setOpen] = useState(false);

  return (
    <>
      <TableRow sx={{ "& > *": { borderBottom: "unset" } }}>
        <TableCell>
          <IconButton size="small" onClick={() => setOpen(!open)}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {new Date(row.created_at).toLocaleString()}
        </TableCell>
        <TableCell>
          {row.user_agent.length > 50 ? `${row.user_agent.substring(0, 50)}...` : row.user_agent}
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Детали
              </Typography>
              <Stack spacing={1}>
                <Box display="flex" justifyContent="space-between">
                  <Typography variant="body2" color="text.secondary">IP:</Typography>
                  <Typography variant="body2">{row.ip || "Неизвестен"}</Typography>
                </Box>
                <Box display="flex" justifyContent="space-between">
                  <Typography variant="body2" color="text.secondary">User Agent:</Typography>
                  <Typography variant="body2">{row.user_agent}</Typography>
                </Box>
              </Stack>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </>
  );
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

    const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
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
            <TableContainer>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell />
                    <TableCell>Дата и время</TableCell>
                    <TableCell>User Agent</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {analytics.clicks.map((click) => (
                    <Row key={click.id} row={click} />
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
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