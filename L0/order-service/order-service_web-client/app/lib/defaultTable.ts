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
  // muiTablePaperProps: {
  //   sx: {
  //     backgroundColor: "background.paper",
  //     boxShadow: "none",
  //   },
  // },
  // muiTableContainerProps: {
  //   sx: {
  //     backgroundColor: "background.paper",
  //   },
  // },
  // muiTableHeadCellProps: {
  //   sx: {
  //     fontSize: "0.5rem",
  //     height: "10px",
  //     backgroundColor: "background.paper",
  //   },
  // },
  // muiBottomToolbarProps:{
  //    sx: {
  //     backgroundColor: "background.paper",
  //   },
  // },
  // muiTopToolbarProps:{
  //   sx: {
  //     backgroundColor: "background.paper",
  //   },
  // },
  // muiTableBodyCellProps: {
  //   sx: {
  //     fontSize: "0.5rem",
  //     minHeight: "10px",
  //     height: "20px",
  //     paddingBottom: "0",
  //     paddingTop: "0",
  //     backgroundColor: "background.paper",
  //   },
  // },
  // muiTopToolbarProps: {
  //   sx: {
  //     "& .MuiButton-root": {
  //       fontSize: "0.75rem",
  //       padding: "4px 8px",
  //       minWidth: "unset",
  //     },
  //   },
  // },
  muiSelectCheckboxProps: {
    size: "small",
  },
  muiSelectAllCheckboxProps: {
    size: "small",
  },
})
