const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL

export const uploadImage = async (file: File): Promise<{ id: string }> => {
  const formData = new FormData()
  formData.append("file", file)
  const res = await fetch(`${API_BASE_URL}/upload`, {
    method: "POST",
    body: formData,
  })
  if (!res.ok) throw new Error("Upload failed")
  return res.json()
}

export const getImageStatus = async (id: string): Promise<Image> => {
  const res = await fetch(`${API_BASE_URL}/image/${id}`)
  if (!res.ok) throw new Error("Failed to fetch image status")
  return res.json()
}

export const deleteImage = async (id: string): Promise<void> => {
  const res = await fetch(`${API_BASE_URL}/image/{id}`, {
    method: "DELETE",
  })
  if (!res.ok) throw new Error("Delete failed")
}
