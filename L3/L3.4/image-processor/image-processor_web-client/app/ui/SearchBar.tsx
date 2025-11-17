"use client"

import { useState } from "react"
import {
  Paper,
  TextField,
  Button,
  Stack,
  Alert,
  InputAdornment,
  IconButton,
} from "@mui/material"
import SearchIcon from "@mui/icons-material/Search"
import { useCommentContext } from "@/app/lib/CommentContext"
import { Comment } from "@/app/lib/types"

export default function SearchBar() {
  const [query, setQuery] = useState("")
  const [error, setError] = useState<string | null>(null)
  const { setComments, fetchComments } = useCommentContext()

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!query.trim()) {
      setError("Введите текст для поиска")
      return
    }

    const apiBase = process.env.NEXT_PUBLIC_API_URL

    try {
      const response = await fetch(
        `${apiBase}/comments/search?q=${encodeURIComponent(query)}`
      )
      if (response.ok) {
        const data = (await response.json()) as Comment[]
        setComments(data)
      } else {
        setError("Ошибка при поиске")
      }
    } catch (err) {
      setError("Ошибка сети при поиске")
    }
  }

  const handleIconClick = () => {
    if (query.trim()) {
      handleSearch({ preventDefault: () => {} } as React.FormEvent)
    }
  }

  const handleReset = () => {
    setQuery("")
    fetchComments()
  }

  return (
    <Paper
      component="form"
      onSubmit={handleSearch}
      sx={{ p: 2, display: "flex", gap: 2, alignItems: "center" }}
    >
      <TextField
        fullWidth
        value={query || ""}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Поиск комментариев..."
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={handleIconClick} edge="end">
                <SearchIcon />
              </IconButton>
            </InputAdornment>
          ),
        }}
      />
      <Button
        onClick={handleReset}
        variant="outlined"
        color="secondary"
        disabled={!query}
      >
        Сброс
      </Button>
      {error && (
        <Alert severity="error" sx={{ mt: 1, width: "100%" }}>
          {error}
        </Alert>
      )}
    </Paper>
  )
}
