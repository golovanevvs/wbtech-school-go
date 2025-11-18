export interface Image {
  id: number
  status: "uploading" | "processing" | "completed" | "failed"
  original_path?: string
  processed_url?: string
  created_at: string
  operations?: {
    resize?: boolean
    thumbnail?: boolean
    watermark?: boolean
  }
}
