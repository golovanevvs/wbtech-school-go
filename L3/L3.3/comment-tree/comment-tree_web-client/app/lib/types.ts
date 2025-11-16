export interface Comment {
  id: number
  parent_id: number | null
  text: string
  created_at: string
  children: Comment[]
}
