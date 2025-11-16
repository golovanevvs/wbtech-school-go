"use client";

import { useState, useEffect } from "react";
import { Stack } from "@mui/material";
import CommentItem from "./CommentItem";
import { Comment } from "@/app/lib/types";

export default function CommentTree() {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchComments = async () => {
      const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      try {
        const response = await fetch(`${apiBase}/comments`);
        if (response.ok) {
          const data = await response.json();
          setComments(data);
        }
      } catch (err) {
        console.error("Ошибка загрузки комментариев:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchComments();
  }, []);

  if (loading) return <div>Загрузка...</div>;

  return (
    <Stack spacing={1}>
      {comments.map((comment) => (
        <CommentItem key={comment.id} comment={comment} />
      ))}
    </Stack>
  );
}