interface OrderData {
  order_uid: string
  track_number: string
  entry: string
  locale: string
  internal_signature: string
  customer_id: string
  delivery_service: string
  shardkey: string
  sm_id: number
  date_created: string
  oof_shard: string
  delivery: Delivery
  payment: Payment
  items: Item[]
}

interface Delivery {
  name: string
  phone: string
  zip: string
  city: string
  address: string
  region: string
  email: string
}

interface Payment {
  transaction: string
  request_id: string
  currency: string
  provider: string
  amount: number
  payment_dt: number
  bank: string
  delivery_cost: number
  goods_total: number
  custom_fee: number
}

export interface Item {
  chrt_id: number
  track_number: string
  price: number
  rid: string
  name: string
  sale: number
  size: string
  total_price: number
  nm_id: number
  brand: string
  status: number
}

export interface OrderResponse {
  success: boolean
  data?: OrderData
  error?: string
}

// interface FieldConfig {
//   label: string
//   value: string | number
//   path?: string
//   color?: string
// }

// export interface SectionConfig {
//   title: string
//   condition?: boolean
//   fields: FieldConfig[]
// }

// Если добавить path в конфиг полей, можно сделать автоматическое получение значений:

// const getValueByPath = (obj: any, path: string) =>
//   path.split('.').reduce((acc, part) => acc?.[part], obj)

// Использование:
// value: getValueByPath(orderResponse, field.path) || "-"

// Можно добавить в конфиг поле component для кастомных компонентов:

// {
//   label: "Дата создания",
//   component: (value) => <DateDisplay value={value} />,
//   path: "data.date_created"
// }
