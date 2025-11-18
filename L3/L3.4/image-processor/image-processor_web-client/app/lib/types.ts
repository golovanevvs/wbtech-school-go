export interface Image {
  id: string
  status: "uploading" | "processing" | "completed" | "failed"
  originalUrl?: string
  processedUrl?: string
  createdAt: string
  operations?: {
    resize?: boolean
    thumbnail?: boolean
    watermark?: boolean
  }
}
