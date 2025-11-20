import { Paper } from "@mui/material"
import { ReactNode } from "react"

interface CardProps {
  children: ReactNode
}

export default function Card({ children }: CardProps) {
  return (
    <Paper
      elevation={3}
      sx={{
        maxWidth: "400px",
        mx: "auto",
        p: 3,
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      {children}
    </Paper>
  )
}
