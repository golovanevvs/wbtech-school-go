"use client"

import { Stack } from "@mui/material"
import CreateNoticeForm from "./CreateNoticeForm"
import GetNoticeStatus from "./GetNoticeStatus"
import DeleteNotice from "./DeleteNotice"
import Header from "./Header"

export default function NotifyForm() {
  return (
    <Stack spacing={4}>
      <Header />
      <CreateNoticeForm />
      <GetNoticeStatus />
      <DeleteNotice />
    </Stack>
  )
}
