export interface SalesRecord {
  id: number
  type: "income" | "expense"
  category: string
  date: string
  amount: number
}

export interface SalesRecordFormData {
  type: "income" | "expense"
  category: string
  date: string
  amount: number
}

export interface CreateSalesRecordResponse {
  id: number
}

export interface AnalyticsData {
  sum: number
  avg: number
  count: number
  median: number
  percentile90: number
}

export interface AnalyticsRequest {
  from: string
  to: string
}

export interface SortOptions {
  field: "id" | "type" | "category" | "date" | "amount"
  direction: "asc" | "desc"
}