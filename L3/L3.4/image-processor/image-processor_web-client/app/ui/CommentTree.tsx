"use client"

import { useEffect } from "react"
import { Stack } from "@mui/material"
import CommentItem from "./CommentItem"
import { useCommentContext } from "@/app/lib/CommentContext"

export default function CommentTree() {
  const { comments, fetchComments } = useCommentContext()

  useEffect(() => {
    fetchComments()
  }, [fetchComments])

  return (
    <Stack spacing={1}>
      {comments.map((comment) => (
        <CommentItem key={comment.id} comment={comment} />
      ))}
    </Stack>
  )
}
