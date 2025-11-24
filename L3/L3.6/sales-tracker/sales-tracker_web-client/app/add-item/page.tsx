"use client"

import { Box } from "@mui/material"
import AddItemForm from "../../ui/addItem/AddItemForm"

export default function AddItemPage() {
  return (
    <Box sx={{ maxWidth: 600, mx: "auto", p: 3 }}>
      <AddItemForm />
    </Box>
  )
}