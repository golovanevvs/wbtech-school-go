"use client"

import { Container } from "@mui/material"
// import styles from "./page.module.css"
import NotifyForm from "./ui/NotifyForm"

export default function Home() {
  return (
    <Container
      disableGutters
      sx={{
        minWidth: "200px",
        maxWidth: "1000px",
        bgcolor: "background.default",
        p:2,
      }}
      >
        <NotifyForm/>
    </Container>
  )
}
