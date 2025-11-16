import { useState } from "react";
import { Paper, Typography, Box, IconButton, Stack } from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import DeleteIcon from "@mui/icons-material/Delete";
import CommentForm from "./CommentForm";
import { Comment } from "@/app/lib/types";

interface CommentItemProps {
  comment: Comment;
  depth?: number;
}

export default function CommentItem({ comment, depth = 0 }: CommentItemProps) {
  const [expanded, setExpanded] = useState(true);
  const [showReply, setShowReply] = useState(false);
  const [deleted, setDeleted] = useState(false);

  if (deleted) return null;

  const handleDelete = async () => {
    const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    try {
      const response = await fetch(`${apiBase}/comments/${comment.id}`, {
        method: "DELETE",
      });
      if (response.ok) {
        setDeleted(true);
      }
    } catch (err) {
      alert("Ошибка при удалении");
    }
  };

  return (
    <Paper
      sx={{
        p: 2,
        ml: depth * 2,
        mt: 1,
        position: "relative",
        borderLeft: depth > 0 ? "2px solid #ccc" : "none",
      }}
    >
      <Stack direction="row" justifyContent="space-between" alignItems="flex-start">
        <Box>
          <Typography variant="body2" color="textSecondary">
            ID: {comment.id}, {new Date(comment.created_at).toLocaleString()}
          </Typography>
          <Typography variant="body1" sx={{ mt: 1 }}>
            {comment.text}
          </Typography>
        </Box>
        <Stack direction="row" spacing={1}>
          <IconButton size="small" onClick={() => setShowReply(!showReply)}>
            <ExpandMoreIcon />
          </IconButton>
          <IconButton size="small" color="error" onClick={handleDelete}>
            <DeleteIcon />
          </IconButton>
        </Stack>
      </Stack>

      {showReply && <CommentForm parentId={comment.id} />}

      {comment.children && comment.children.length > 0 && expanded && (
        <Box sx={{ mt: 1 }}>
          {comment.children.map((child) => (
            <CommentItem key={child.id} comment={child} depth={depth + 1} />
          ))}
        </Box>
      )}
    </Paper>
  );
}