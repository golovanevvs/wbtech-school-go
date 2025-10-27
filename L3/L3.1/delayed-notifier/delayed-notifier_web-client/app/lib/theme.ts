import { createTheme } from "@mui/material/styles"
import { Geist } from "next/font/google"

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
})

export const baseTheme = createTheme({
  typography: {
    fontFamily: geist.style.fontFamily,
  },
  components: {
    MuiTypography: {
      styleOverrides: {
        h1: {
          textAlign: "center",
          marginBottom: 20,
          fontSize: 30,
          fontWeight: "bold",
        },
        h2: {
          textAlign: "center",
          marginBottom: 20,
          fontSize: 18,
          fontWeight: "bold",
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          display: "flex",
          flexDirection: "column",
          borderRadius: 12,
          padding: 20,
          marginBottom: 20,
          boxShadow: "0px 4px 12px rgba(0,0,0,0.1)",
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          "& .MuiOutlinedInput-root": {
            borderRadius: 8,
            transition: "border-color 0.2s, box-shadow 0.2s",
            "&:hover fieldset": {
              borderColor: "primary.main",
            },
            "&.Mui-focused fieldset": {
              borderColor: "primary.main",
            },
          },
        },
      },
    },
  },
})

export const lightTheme = createTheme({
  ...baseTheme,
  palette: {
    mode: "light",
    primary: {
      main: "#6a1b9a",
      light: "#9c4dcc",
      dark: "#4a0072",
    },
    secondary: {
      main: "#ff9800",
      light: "#ffc947",
      dark: "#c66900",
    },
    error: {
      main: "#d32f2f",
    },
    success: {
      main: "#2e7d32",
    },
    background: {
      default: "#f9f4fb",
      paper: "#ffffff",
    },
  },
})

export const darkTheme = createTheme({
  ...baseTheme,
  palette: {
    mode: "dark",
    primary: {
      main: "#ce93d8",
      light: "#f3e5f5",
      dark: "#9c4dcc",
    },
    secondary: {
      main: "#ffb74d",
      light: "#ffe97d",
      dark: "#c88719",
    },
    error: {
      main: "#ef5350",
    },
    success: {
      main: "#81c784",
    },
    background: {
      default: "#1a0f1d",
      paper: "#2c1b2f",
    },
  },
})
