import { TextField, TextFieldProps } from "@mui/material"

interface InputProps extends Omit<TextFieldProps, "label"> {
  label: string
}

export default function Input({ label, ...props }: InputProps) {
  return <TextField label={label} fullWidth margin="normal" {...props} />
}
