"use client"

import { Stack } from "@mui/material"
import CreateNoticeForm from "./CreateNoticeForm"
import GetNoticeStatus from "./GetNoticeStatus"
import DeleteNotice from "./DeleteNotice"

export default function NotifyForm() {
  return (
    <Stack spacing={4}>
      <CreateNoticeForm />
      <GetNoticeStatus />
      <DeleteNotice />
    </Stack>
  )
}
