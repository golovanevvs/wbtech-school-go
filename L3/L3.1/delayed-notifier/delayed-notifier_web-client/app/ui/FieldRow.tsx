"use client"

import { Box } from "@mui/material"

interface FieldRowProps {
  label: string
  value: string | number | null
  statusColor?: string
}

export const FieldRow = ({ label, value, statusColor }: FieldRowProps) => {
  return (
    <Box
      sx={{
        display: "flex",
        mb: 1,
        alignItems: "center",
      }}
    >
      <Box
        sx={{
          fontWeight: "bold",
          minWidth: "150px",
          color: "primary.main",
        }}
      >
        {label}:
      </Box>
      <Box
        sx={{
          flex: "1",
          p: "2px 10px",
          borderRadius: "4px",
          fontSize: "12px",
          backgroundColor: "background.default",
          color: `${statusColor}`,
        }}
      >
        {value}
      </Box>
    </Box>
  )
}
