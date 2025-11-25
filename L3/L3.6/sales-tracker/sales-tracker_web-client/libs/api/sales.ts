import { SalesRecord, SalesRecordFormData, AnalyticsData, AnalyticsRequest, SortOptions, CreateSalesRecordResponse, UpdateSalesRecordResponse, DeleteSalesRecordResponse } from "../types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export async function createSalesRecord(data: SalesRecordFormData): Promise<CreateSalesRecordResponse> {
  const response = await fetch(`${API_BASE_URL}/items`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    throw new Error(`Failed to create sales record: ${response.statusText}`)
  }

  return response.json()
}

export async function getSalesRecords(sortOptions?: SortOptions): Promise<SalesRecord[]> {
  const params = new URLSearchParams()
  
  if (sortOptions) {
    params.append("field", sortOptions.field)
    params.append("direction", sortOptions.direction)
  }

  const response = await fetch(`${API_BASE_URL}/items?${params}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to fetch sales records: ${response.statusText}`)
  }

  const data = await response.json()
  return data.records
}

export async function updateSalesRecord(id: number, data: Partial<SalesRecordFormData>): Promise<UpdateSalesRecordResponse> {
  const response = await fetch(`${API_BASE_URL}/items/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    throw new Error(`Failed to update sales record: ${response.statusText}`)
  }

  return response.json()
}

export async function deleteSalesRecord(id: number): Promise<DeleteSalesRecordResponse> {
  const response = await fetch(`${API_BASE_URL}/items/${id}`, {
    method: "DELETE",
  })

  if (!response.ok) {
    throw new Error(`Failed to delete sales record: ${response.statusText}`)
  }

  return response.json()
}

export async function getAnalytics(request: AnalyticsRequest): Promise<AnalyticsData> {
  const params = new URLSearchParams({
    from: request.from,
    to: request.to,
  })

  const response = await fetch(`${API_BASE_URL}/analytics?${params}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to fetch analytics: ${response.statusText}`)
  }

  return response.json()
}

export async function downloadCSV(from?: string, to?: string): Promise<Blob> {
  const params = new URLSearchParams()
  
  if (from) params.append("from", from)
  if (to) params.append("to", to)

  const response = await fetch(`${API_BASE_URL}/items/export?${params}`, {
    method: "GET",
  })

  if (!response.ok) {
    throw new Error(`Failed to download CSV: ${response.statusText}`)
  }

  return response.blob()
}