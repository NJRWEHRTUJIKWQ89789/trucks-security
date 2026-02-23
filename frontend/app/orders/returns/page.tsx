"use client";

import { useEffect, useState } from "react";
import { Package, Clock, CheckCircle, DollarSign, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Order {
  id: string;
  orderNumber: string;
  customerName: string;
  status: string;
  totalAmount: number;
  returnReason: string;
  createdAt: string;
  updatedAt: string;
}

const statusBadge: Record<string, string> = {
  returned: "bg-orange-100 text-orange-700",
  pending: "bg-amber-100 text-amber-700",
  processing: "bg-blue-100 text-blue-700",
  approved: "bg-green-100 text-green-700",
  refunded: "bg-purple-100 text-purple-700",
  rejected: "bg-red-100 text-red-700",
};

export default function ReturnsPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ returnOrders: { items: Order[]; totalCount: number } }>(`{
      returnOrders(page: 1, perPage: 50) {
        items {
          id orderNumber customerName status totalAmount returnReason createdAt updatedAt
        }
        totalCount
      }
    }`)
      .then((d) => {
        setOrders(d.returnOrders.items);
        setTotalCount(d.returnOrders.totalCount);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  const totalRefundValue = orders.reduce((sum, o) => sum + (o.totalAmount || 0), 0);

  const summaryCards = [
    { label: "Total Returns", value: totalCount, icon: Package, bg: "bg-gray-100", text: "text-gray-600" },
    { label: "Returned", value: orders.filter((o) => o.status === "returned").length, icon: Clock, bg: "bg-amber-100", text: "text-amber-600" },
    { label: "Total Items", value: orders.length, icon: CheckCircle, bg: "bg-green-100", text: "text-green-600" },
    {
      label: "Total Value",
      value: `$${totalRefundValue.toLocaleString("en-US", { minimumFractionDigits: 0 })}`,
      icon: DollarSign,
      bg: "bg-blue-100",
      text: "text-blue-600",
    },
  ];

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Returns Management</h1>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {summaryCards.map((card) => (
          <div key={card.label} className="rounded-lg border border-gray-200 bg-white p-5">
            <div className="flex items-center gap-3">
              <div className={`rounded-lg p-2.5 ${card.bg}`}>
                <card.icon className={`h-5 w-5 ${card.text}`} />
              </div>
              <div>
                <p className="text-sm text-gray-500">{card.label}</p>
                <p className="text-2xl font-bold text-gray-900">{card.value}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="overflow-hidden rounded-lg border border-gray-200 bg-white">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Order ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Customer</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Reason</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Amount</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.map((item) => (
              <tr key={item.id} className="hover:bg-gray-50">
                <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-blue-600">{item.orderNumber}</td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">{item.customerName || "Unknown"}</td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">{item.returnReason || "N/A"}</td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                  {new Date(item.updatedAt || item.createdAt).toLocaleDateString()}
                </td>
                <td className="whitespace-nowrap px-6 py-4">
                  <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusBadge[item.status] || "bg-gray-100 text-gray-700"}`}>
                    {item.status}
                  </span>
                </td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                  {item.totalAmount > 0
                    ? `$${item.totalAmount.toLocaleString("en-US", { minimumFractionDigits: 2 })}`
                    : "\u2014"}
                </td>
              </tr>
            ))}
            {orders.length === 0 && (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-gray-400">
                  No returns found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
