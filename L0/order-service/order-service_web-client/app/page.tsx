"use client"

import { Container, Box, Paper } from "@mui/material"
import Typography from "@mui/material/Typography"
import HookForm from "@/app/ui/HookForm"
import { FieldRow } from "@/app/ui/FieldRow"
import { useState } from "react"
import { OrderResponse } from "@/app/lib/types"
import { getOrderSections } from "@/app/lib/data"
import { ItemsTable } from "@/app/ui/ItemsTable"

export default function Home() {
  const [orderResponse, setOrderResponse] = useState<OrderResponse | null>(null)

  const handleOrderDataReceived = (response: OrderResponse) => {
    setOrderResponse(response)
  }

  const dataSections = getOrderSections(orderResponse)

  return (
    <Container
      disableGutters
      sx={{
        width: "1800",
        bgcolor: "background.default",
      }}
    >
      <Container sx={{ width: "550px" }}>
        <Typography variant="h1" component="h1" sx={{ color: "primary.dark" }}>
          Получение данных заказа
        </Typography>
        <HookForm onOrderDataReceived={handleOrderDataReceived} />
      </Container>

      {orderResponse?.success && (
        <Container
          maxWidth={false}
          sx={{ display: "flex", flexDirection: "column", minWidth: "550px" }}
        >
          <Box
            sx={{
              display: "flex",
              flexDirection: "row",
              justifyContent: "space-between",
              gap: 2,
              width: "100%",
              flexWrap: "wrap",
            }}
          >
            {dataSections.map(
              (section, index) =>
                section.condition !== false && (
                  <Paper
                    key={index}
                    sx={{
                      flex: "1 1 30%",
                      minWidth: "350px",
                      maxWidth: "100%",
                    }}
                  >
                    <Typography
                      variant="h2"
                      component="h2"
                      sx={{ color: "primary.dark" }}
                    >
                      {section.title}
                    </Typography>

                    {section.fields.map((field, fieldIndex) => (
                      <FieldRow
                        key={fieldIndex}
                        label={field.label}
                        value={field.value}
                        statusColor="primary.main"
                      />
                    ))}
                  </Paper>
                )
            )}
          </Box>
          <Paper>
            <Typography
              variant="h2"
              component="h2"
              sx={{ color: "primary.dark" }}
            >
              Товары
            </Typography>
            <ItemsTable orderResponse={orderResponse} />
          </Paper>
        </Container>
      )}
    </Container>
  )
}
