"use client"

import { Container, Stack } from "@mui/material"
import Header from "./ui/Header"
import SearchBar from "@/app/ui/SearchBar"
import CommentTree from "@/app/ui/CommentTree"
import CommentForm from "@/app/ui/CommentForm"

export default function Home() {
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
        <SearchBar />
        <CommentForm parentId={null} />
        <CommentTree />
      </Stack>
    </Container>
  )
}
