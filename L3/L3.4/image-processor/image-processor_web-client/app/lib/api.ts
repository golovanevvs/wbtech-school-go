const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL

import { Image } from "./types"

interface ProcessOptions {
  resize: boolean
  thumbnail: boolean
  watermark: boolean
}

export const uploadImage = async (
  file: File,
  options: ProcessOptions
): Promise<{ id: string }> => {
  const formData = new FormData()
  formData.append("file", file)
  formData.append("options", JSON.stringify(options))

  const res = await fetch(`${API_BASE_URL}/upload`, {
    method: "POST",
    body: formData,
  })
  if (!res.ok) throw new Error("Upload failed")
  return res.json()
}

export const getImageStatus = async (id: string): Promise<Image> => {
  console.log("getImageStatus called for id:", id)

  const res = await fetch(`${API_BASE_URL}/image/${id}`)

  if (!res.ok) {
    const errorText = await res.text()
    console.error("API Error:", errorText)
    throw new Error(`Failed to fetch image status: ${errorText}`)
  }

  const data = await res.json() // ✅ Читаем один раз
  console.log("getImageStatus response:", res.status, data)

  return data
}

export const getAllImages = async (): Promise<Image[]> => {
  console.log("getAllImages called")

  const res = await fetch(`${API_BASE_URL}/images`)

  if (!res.ok) {
    const errorText = await res.text()
    console.error("API Error:", errorText)
    throw new Error(`Failed to fetch all images: ${errorText}`)
  }

  const data = await res.json()
  console.log("getAllImages response:", res.status, data)

  return data
}

export const deleteImage = async (id: string): Promise<void> => {
  const res = await fetch(`${API_BASE_URL}/image/${id}`, {
    method: "DELETE",
  })
  if (!res.ok) throw new Error("Delete failed")
}