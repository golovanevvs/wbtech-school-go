"use client"

import { useState, useEffect, useCallback } from "react"
import { Box, Typography, Alert, FormControl, InputLabel, Select, MenuItem, Button, IconButton } from "@mui/material"
import { MaterialReactTable, type MRT_ColumnDef } from "material-react-table"
import { Edit, Delete } from "@mui/icons-material"
import { SalesRecord, SortOptions } from "../../libs/types"
import { getSalesRecords, updateSalesRecord, deleteSalesRecord } from "../../libs/api/sales"
import EditItemForm from "./EditItemForm"

export default function ItemsList() {
  const [records, setRecords] = useState<SalesRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [editingRecord, setEditingRecord] = useState<SalesRecord | null>(null)
  const [sortField, setSortField] = useState<"id" | "type" | "category" | "date" | "amount">("id")
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("desc")

  const fetchRecords = useCallback(async () => {
    try {
      setLoading(true)
      const sortOptions: SortOptions = {
        field: sortField,
        direction: sortDirection,
      }
      const data = await getSalesRecords(sortOptions)
      setRecords(data)
      setError("")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка при загрузке записей")
    } finally {
      setLoading(false)
    }
  }, [sortField, sortDirection])

  useEffect(() => {
    fetchRecords()
  }, [fetchRecords])

  const handleSort = async () => {
    const newDirection = sortDirection === "asc" ? "desc" : "asc"
    setSortDirection(newDirection)
  }

  const handleEdit = (record: SalesRecord) => {
    setEditingRecord(record)
  }

  const handleDelete = async (id: number) => {
    if (window.confirm("Вы уверены, что хотите удалить эту запись?")) {
      try {
        await deleteSalesRecord(id)
        await fetchRecords()
      } catch (err) {
        setError(err instanceof Error ? err.message : "Ошибка при удалении записи")
      }
    }
  }

  const handleEditSave = async (id: number, data: Partial<SalesRecord>) => {
    try {
      await updateSalesRecord(id, data)
      setEditingRecord(null)
      await fetchRecords()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка при обновлении записи")
    }
  }

  const handleEditCancel = () => {
    setEditingRecord(null)
  }

  const columns: MRT_ColumnDef<SalesRecord>[] = [
    {
      accessorKey: "id",
      header: "ID",
      size: 80,
    },
    {
      accessorKey: "type",
      header: "Тип",
      size: 100,
      Cell: ({ cell }) => (
        <Box
          component="span"
          sx={{
            px: 1,
            py: 0.5,
            borderRadius: 1,
            fontSize: "0.75rem",
            fontWeight: "bold",
            bgcolor: cell.getValue<string>() === "income" ? "#4caf50" : "#f44336",
            color: "white",
          }}
        >
          {cell.getValue<string>() === "income" ? "Доход" : "Расход"}
        </Box>
      ),
    },
    {
      accessorKey: "category",
      header: "Категория",
      size: 150,
    },
    {
      accessorKey: "date",
      header: "Дата",
      size: 180,
      Cell: ({ cell }) => {
        const date = new Date(cell.getValue<string>())
        return date.toLocaleDateString("ru-RU")
      },
    },
    {
      accessorKey: "amount",
      header: "Сумма",
      size: 120,
      Cell: ({ cell }) => {
        const amount = cell.getValue<number>()
        return `${amount.toLocaleString("ru-RU", { minimumFractionDigits: 2 })} ₽`
      },
    },
    {
      id: "actions",
      header: "Действия",
      size: 100,
      Cell: ({ row }) => (
        <Box sx={{ display: "flex", gap: 1 }}>
          <IconButton
            color="primary"
            onClick={() => handleEdit(row.original)}
            size="small"
          >
            <Edit />
          </IconButton>
          <IconButton
            color="error"
            onClick={() => handleDelete(row.original.id)}
            size="small"
          >
            <Delete />
          </IconButton>
        </Box>
      ),
    },
  ]

  return (
    <Box sx={{ width: "100%" }}>
      <Typography variant="h4" component="h1" gutterBottom align="center">
        Список записей
      </Typography>

      <Box sx={{ mb: 2, display: "flex", gap: 2, alignItems: "center" }}>
        <FormControl sx={{ minWidth: 120 }}>
          <InputLabel>Поле сортировки</InputLabel>
          <Select
            value={sortField}
            label="Поле сортировки"
            onChange={(e) => setSortField(e.target.value as "id" | "type" | "category" | "date" | "amount")}
          >
            <MenuItem value="id">ID</MenuItem>
            <MenuItem value="type">Тип</MenuItem>
            <MenuItem value="category">Категория</MenuItem>
            <MenuItem value="date">Дата</MenuItem>
            <MenuItem value="amount">Сумма</MenuItem>
          </Select>
        </FormControl>

        <Button variant="outlined" onClick={handleSort}>
          {sortDirection === "asc" ? "По возрастанию" : "По убыванию"}
        </Button>

        <Button variant="contained" onClick={fetchRecords}>
          Обновить
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <MaterialReactTable
        columns={columns}
        data={records}
        state={{ isLoading: loading }}
        initialState={{
          pagination: {
            pageIndex: 0,
            pageSize: 10,
          },
        }}
        muiTableContainerProps={{
          sx: {
            width: "100%",
            minHeight: "400px",
          },
        }}
        muiTableProps={{
          sx: {
            width: "100%",
          },
        }}
      />

      {editingRecord && (
        <EditItemForm
          record={editingRecord}
          onSave={handleEditSave}
          onCancel={handleEditCancel}
        />
      )}
    </Box>
  )
}