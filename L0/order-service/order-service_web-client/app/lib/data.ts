import { OrderResponse } from "@/app/lib/types"

export const getOrderSections = (orderResponse: OrderResponse | null) => [
  {
    title: "Основная информация о заказе",
    condition: orderResponse?.success,
    fields: [
      {
        label: "Order UID",
        value: orderResponse?.data?.order_uid || "-",
      },
      {
        label: "Track Number",
        value: orderResponse?.data?.track_number || "-",
      },
      {
        label: "Entry",
        value: orderResponse?.data?.entry || "-",
      },
      {
        label: "Locale",
        value: orderResponse?.data?.locale || "-",
      },
      {
        label: "Internal Signature",
        value: orderResponse?.data?.internal_signature || "-",
      },
      {
        label: "Customer ID",
        value: orderResponse?.data?.customer_id || "-",
      },
      {
        label: "Delivery Service",
        value: orderResponse?.data?.delivery_service || "-",
      },
      {
        label: "Shardkey",
        value: orderResponse?.data?.shardkey || "-",
      },
      {
        label: "Sm ID",
        value: orderResponse?.data?.sm_id || "-",
      },
      {
        label: "Date Created",
        value: orderResponse?.data?.date_created
          ? new Date(orderResponse.data.date_created).toLocaleString("ru-RU")
          : "-",
      },
      {
        label: "OOF Shard",
        value: orderResponse?.data?.oof_shard || "-",
      },
    ],
  },
  {
    title: "Сведения о доставке",
    condition: orderResponse?.success,
    fields: [
      {
        label: "Name",
        value: orderResponse?.data?.delivery.name || "-",
      },
      {
        label: "Phone",
        value: orderResponse?.data?.delivery.phone || "-",
      },
      {
        label: "ZIP",
        value: orderResponse?.data?.delivery.zip || "-",
      },
      {
        label: "City",
        value: orderResponse?.data?.delivery.city || "-",
      },
      {
        label: "Address",
        value: orderResponse?.data?.delivery.address || "-",
      },
      {
        label: "Region",
        value: orderResponse?.data?.delivery.region || "-",
      },
      {
        label: "Email",
        value: orderResponse?.data?.delivery.email || "-",
      },
    ],
  },
  {
    title: "Сведения об оплате",
    condition: orderResponse?.success,
    fields: [
      {
        label: "Transaction",
        value: orderResponse?.data?.payment.transaction || "-",
      },
      {
        label: "Request ID",
        value: orderResponse?.data?.payment.request_id || "-",
      },
      {
        label: "Currency",
        value: orderResponse?.data?.payment.currency || "-",
      },
      {
        label: "Provider",
        value: orderResponse?.data?.payment.provider || "-",
      },
      {
        label: "Amount",
        value: orderResponse?.data?.payment.amount || "-",
      },
      {
        label: "Payment DT",
        value: orderResponse?.data?.payment.payment_dt || "-",
      },
      {
        label: "Bank",
        value: orderResponse?.data?.payment.bank || "-",
      },
      {
        label: "Delivery Cost",
        value: orderResponse?.data?.payment.delivery_cost || "-",
      },
      {
        label: "Goods Total",
        value: orderResponse?.data?.payment.goods_total || "-",
      },
      {
        label: "Custom Fee",
        value: orderResponse?.data?.payment.custom_fee || "-",
      },
    ],
  },
]
