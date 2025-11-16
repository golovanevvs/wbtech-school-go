"use client";

import { Box } from "@mui/material";
import { ReactNode } from "react";

interface FieldRowProps {
  label: string;
  value: string | number | null | ReactNode;
  statusColor?: string;
}

export const FieldRow = ({ label, value, statusColor }: FieldRowProps) => {
  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: { xs: "column", sm: "row" },
        mb: 1,
        alignItems: "flex-start",
      }}
    >
      <Box
        sx={{
          fontWeight: "bold",
          minWidth: "100px",
          color: "primary.main",
          mb: { xs: 0.5, sm: 0 },
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
  );
};