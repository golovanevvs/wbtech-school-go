// Типы для работы с товарами склада

export interface Item {
  id: number
  name: string
  price: number
  quantity: number
  created_at: string
  updated_at: string
}

// Тип для ответа API при получении одного товара
export interface ItemResponse {
  item: Item
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
  changes?: string
  created_at: string
}

export interface ItemHistoryResponse {
  history: ItemAction[]
}
