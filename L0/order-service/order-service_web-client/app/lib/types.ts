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