"use client"

import { Box} from "@mui/material"
import ItemsList from "../ui/items/ItemsList"

export default function Home() {
 

   return (
    <Box sx={{ p: 3 }}>
      {/* Main content - Items List */}
      <ItemsList />
    </Box>
  )
}