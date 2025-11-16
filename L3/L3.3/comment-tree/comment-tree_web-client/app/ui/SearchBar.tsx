"use client";

import { useState } from "react";
import { Paper, TextField, Button, Stack } from "@mui/material";

export default function SearchBar() {
  const [query, setQuery] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Search query:", query);
    // Здесь будет вызов API
  };

  return (
    <Paper
      component="form"
      onSubmit={handleSearch}
      sx={{ p: 2, display: "flex", gap: 2, alignItems: "center" }}
    >
      <TextField
        fullWidth
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Поиск комментариев..."
      />
      <Button type="submit" variant="contained">
        Найти
      </Button>
    </Paper>
  );
}