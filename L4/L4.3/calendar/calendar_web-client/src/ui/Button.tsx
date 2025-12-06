import {
  Button as MuiButton,
  ButtonProps as MuiButtonProps,
} from "@mui/material"

interface ButtonProps extends MuiButtonProps {
  children: React.ReactNode
}

export default function Button({ children, ...props }: ButtonProps) {
  return (
    <MuiButton variant="contained" fullWidth {...props}>
      {children}
    </MuiButton>
  )
}