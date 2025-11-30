import apiClient from "./client"
import { Item, ItemsResponse, ItemHistoryResponse, ItemResponse } from "../types/items"

// API для работы с товарами склада
export const itemsAPI = {
  /**
   * Получение списка всех товаров
   * @returns Список товаров
   */
  async getItems(): Promise<ItemsResponse> {
    try {
      const response = await apiClient.get<ItemsResponse>("/items")
      return response
    } catch (error) {
      console.error("Failed to fetch items:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось загрузить список товаров")
    }
  },

  /**
   * Получение одного товара по ID
   * @param itemId ID товара
   * @returns Товар
   */
  async getItem(itemId: number): Promise<ItemResponse> {
    try {
      const response = await apiClient.get<ItemResponse>(`/items/${itemId}`)
      return response
    } catch (error) {
      console.error("Failed to fetch item:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось загрузить товар")
    }
  },

  /**
   * Получение истории действий с товаром
   * @param itemId ID товара
   * @returns История действий
   */
  async getItemHistory(itemId: number): Promise<ItemHistoryResponse> {
    try {
      const response = await apiClient.get<ItemHistoryResponse>(`/items/${itemId}/history`)
      return response
    } catch (error) {
      console.error("Failed to fetch item history:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось загрузить историю товара")
    }
  },

  /**
   * Создание нового товара
   * @param itemData Данные нового товара
   * @returns Созданный товар
   */
  async createItem(itemData: Omit<Item, 'id' | 'created_at' | 'updated_at'>): Promise<Item> {
    try {
      const response = await apiClient.post<Item>("/items", itemData)
      return response
    } catch (error) {
      console.error("Failed to create item:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось создать товар")
    }
  },

  /**
   * Обновление товара
   * @param itemId ID товара
   * @param itemData Обновленные данные товара
   * @returns Обновленный товар
   */
  async updateItem(itemId: number, itemData: Partial<Omit<Item, 'id' | 'created_at' | 'updated_at'>>): Promise<Item> {
    try {
      const response = await apiClient.put<Item>(`/items/${itemId}`, itemData)
      return response
    } catch (error) {
      console.error("Failed to update item:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось обновить товар")
    }
  },

  /**
   * Удаление товара
   * @param itemId ID товара
   */
  async deleteItem(itemId: number): Promise<void> {
    try {
      await apiClient.delete(`/items/${itemId}`)
    } catch (error) {
      console.error("Failed to delete item:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось удалить товар")
    }
  },

  /**
   * Экспорт истории товара в CSV
   * @param itemId ID товара
   * @returns Blob с CSV данными
   */
  async exportItemHistoryCSV(itemId: number): Promise<Blob> {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/items/${itemId}/history/export`, {
        method: 'GET',
        credentials: 'include',
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.blob()
    } catch (error) {
      console.error("Failed to export CSV:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось экспортировать CSV")
    }
  }
}

export type { ItemsResponse, ItemHistoryResponse }