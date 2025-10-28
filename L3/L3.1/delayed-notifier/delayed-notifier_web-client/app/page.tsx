"use client"

import { Container } from "@mui/material"
// import styles from "./page.module.css"
import NotifyForm from "./ui/NotifyForm"

export default function Home() {
  return (
    <Container
      disableGutters
      sx={{
        width: "1800",
        bgcolor: "background.default",
      }}
      >
        <NotifyForm/>
    </Container>
  )
}
