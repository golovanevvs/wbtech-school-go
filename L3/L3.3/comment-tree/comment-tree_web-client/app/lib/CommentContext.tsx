"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { Comment } from "./types";

type CommentContextType = {
  comments: Comment[];
  setComments: (comments: Comment[]) => void;
  fetchComments: () => Promise<void>;
};

const CommentContext = createContext<CommentContextType | undefined>(undefined);

export function CommentProvider({ children }: { children: ReactNode }) {
  const [comments, setComments] = useState<Comment[]>([]);

  const fetchComments = async () => {
    const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    try {
      const response = await fetch(`${apiBase}/comments`);
      if (response.ok) {
        const data = await response.json();
        setComments(data);
      } else {
        console.error("Ошибка при загрузке комментариев:", response.status);
      }
    } catch (err) {
      console.error("Ошибка сети:", err);
    }
  };

  return (
    <CommentContext.Provider value={{ comments, setComments, fetchComments }}>
      {children}
    </CommentContext.Provider>
  );
}

export function useCommentContext() {
  const context = useContext(CommentContext);
  if (!context) {
    throw new Error("useCommentContext must be used within a CommentProvider");
  }
  return context;
}