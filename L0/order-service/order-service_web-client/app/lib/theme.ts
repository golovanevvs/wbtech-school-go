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

// Светлая тема
export const lightTheme = createTheme({
  ...baseTheme,
  palette: {
    mode: "light",
    primary: {
      main: "#6a1b9a",   // Фиолетовый, но мягкий
      light: "#9c4dcc",  // Светлый акцент
      dark: "#4a0072",   // Глубокий фиолет
    },
    secondary: {
      main: "#ff9800",   // Тёплый оранжевый
      light: "#ffc947",  // Мягкий жёлто-оранжевый
      dark: "#c66900",   // Глубокий оранжевый
    },
    error: {
      main: "#d32f2f",   // Красный, но не кислотный
    },
    success: {
      main: "#2e7d32",   // Зелёный в сторону изумрудного
    },
    background: {
      default: "#f9f4fb", // Очень светлый фиолетово-серый фон
      paper: "#ffffff",   // Белая карточка
    },
  },
})

// Тёмная тема
export const darkTheme = createTheme({
  ...baseTheme,
  palette: {
    mode: "dark",
    primary: {
      main: "#ce93d8",   // Пастельный фиолетовый
      light: "#f3e5f5",  // Очень светлый фиолет
      dark: "#9c4dcc",   // Более насыщенный
    },
    secondary: {
      main: "#ffb74d",   // Мягкий тёплый
      light: "#ffe97d",  // Светлый золотистый
      dark: "#c88719",   // Тёмный золотой
    },
    error: {
      main: "#ef5350",   // Красный без излишней агрессии
    },
    success: {
      main: "#81c784",   // Мягкий зелёный
    },
    background: {
      default: "#1a0f1d", // Глубокий фиолетово-серый
      paper: "#2c1b2f",   // Карточка чуть светлее фона
    },
  },
})
