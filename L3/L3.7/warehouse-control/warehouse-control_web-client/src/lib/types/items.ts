// Типы для работы с товарами склада

export interface Item {
  id: number
  name: string
  price: number
  quantity: number
  created_at: string
  updated_at: string
}

export interface ItemsResponse {
  items: Item[]
}

export interface ItemAction {
  id: number
  item_id: number
  action_type: "create" | "update" | "delete"
  user_id: number
  user_name: string
  changes?: Record<string, unknown>
  created_at: string
}

export interface ItemHistoryResponse {
  history: ItemAction[]
}