"use client"

import { useMemo } from "react"
import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from "material-react-table"

type Order = {
  id: number
  name: string
  comm: string
}

const data: Order[] = [
  {
    id: 1,
    name: "Bosch",
    comm: "Paid",
  },
  {
    id: 2,
    name: "Sheba",
    comm: "Paid",
  },
  {
    id: 3,
    name: "CRISP",
    comm: "Trunch",
  },
]

const Table = () => {
  const columns = useMemo<MRT_ColumnDef<Order>[]>(
    () => [
      {
        accessorKey: "id",
        header: "ID number",
        size: 50,
      },
      {
        accessorKey: "name",
        header: "Name",
        size: 100,
      },
      {
        accessorKey: "comm",
        header: "Comment",
        size: 200,
      },
    ],
    []
  )

  const table = useMaterialReactTable({
    columns,
    data,
  })

  return <MaterialReactTable table={table} />
}

export default Table
