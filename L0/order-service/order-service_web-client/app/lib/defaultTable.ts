import { type MRT_RowData, type MRT_TableOptions } from "material-react-table"

export const getDefaultMRTOptions = <TData extends MRT_RowData>(): Partial<
  MRT_TableOptions<TData>
> => ({
  muiTablePaperProps: {
    sx: {
      boxShadow: "none",
    },
  },
  muiTableHeadCellProps: {
    align: "center",
    sx: {
      padding: "0 0",
      "& .Mui-TableHeadCell-Content": {
        gap: "0px",
      },
    },
  },
  muiTableBodyCellProps: {
    sx: {
      textAlign: "center",
    },
  },
  muiSelectCheckboxProps: {
    size: "small",
  },
  muiSelectAllCheckboxProps: {
    size: "small",
  },
})
