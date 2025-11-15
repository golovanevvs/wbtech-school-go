export interface ShortURL {
  id: number
  original: string
  short: string
  custom: boolean
  created_at: string
}

export interface ClickEvent {
  id: number
  short: string
  user_agent: string
  ip: string
  created_at: string
}

export interface AnalyticsData {
  total_clicks: number
  clicks: ClickEvent[]
}
