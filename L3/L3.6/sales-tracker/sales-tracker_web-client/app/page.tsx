"use client"

import { Box} from "@mui/material"
import ItemsList from "../ui/items/ItemsList"

export default function Home() {
 

   return (
    <Box sx={{ width: "100%" }}>
      {/* Main content - Items List */}
      <ItemsList />
    </Box>
  )
}