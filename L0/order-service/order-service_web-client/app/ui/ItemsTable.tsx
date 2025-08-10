"use client"

import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from "material-react-table"
import { MRT_Localization_RU } from "material-react-table/locales/ru"
import { getDefaultMRTOptions } from "@/app/lib/defaultTable"
import { OrderResponse, Item } from "@/app/lib/types"

interface ItemsTableProps {
  orderResponse: OrderResponse | null
}

export const ItemsTable = ({ orderResponse }: ItemsTableProps) => {
  const items: Item[] = orderResponse?.data?.items || []

  const columns: MRT_ColumnDef<Item>[] = [
    {
      accessorKey: "chrt_id",
      header: "chrt_id",
      size: 30,
    },
    {
      accessorKey: "track_number",
      header: "track_number",
      size: 100,
    },
    {
      accessorKey: "price",
      header: "price",
      size: 20,
    },
    {
      accessorKey: "rid",
      header: "rid",
      size: 60,
    },
    {
      accessorKey: "name",
      header: "name",
      size: 60,
    },
    {
      accessorKey: "sale",
      header: "sale",
      size: 60,
    },
    {
      accessorKey: "size",
      header: "size",
      size: 20,
    },
    {
      accessorKey: "total_price",
      header: "total_price",
      size: 60,
    },
    {
      accessorKey: "nm_id",
      header: "nm_id",
      size: 60,
    },
    {
      accessorKey: "brand",
      header: "brand",
      size: 60,
    },
    {
      accessorKey: "status",
      header: "status",
      size: 60,
    },
  ]

  const defaultMRTOptions = getDefaultMRTOptions<Item>()
  const table = useMaterialReactTable({
    ...defaultMRTOptions,
    columns,
    data: items,
    enableGlobalFilter: true,
    localization: MRT_Localization_RU,
    initialState: {
      ...defaultMRTOptions.initialState,
      showColumnFilters: false,
    },
  })

  return <MaterialReactTable table={table} />
}
