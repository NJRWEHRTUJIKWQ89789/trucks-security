"use client";

import { useEffect, useState } from "react";
import { XCircle, DollarSign, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Order {
  id: string;
  orderNumber: string;
  customerName: string;
  status: string;
  totalAmount: number;
  cancellationReason: string;
  createdAt: string;
  updatedAt: string;
}

const statusBadge: Record<string, string> = {
  cancelled: "bg-red-100 text-red-700",
  refunded: "bg-green-100 text-green-700",
  processing: "bg-amber-100 text-amber-700",
  denied: "bg-red-100 text-red-700",
};

export default function CancellationsPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ cancelledOrders: { items: Order[]; totalCount: number } }>(`{
      cancelledOrders(page: 1, perPage: 50) {
        items {
          id orderNumber customerName status totalAmount cancellationReason createdAt updatedAt
        }
        totalCount
      }
    }`)
      .then((d) => {
        setOrders(d.cancelledOrders.items);
        setTotalCount(d.cancelledOrders.totalCount);
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

  const valueLost = orders.reduce((sum, o) => sum + (o.totalAmount || 0), 0);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Cancellations</h1>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="rounded-lg border border-gray-200 bg-white p-5">
          <div className="flex items-center gap-3">
            <div className="rounded-lg bg-red-100 p-2.5">
              <XCircle className="h-5 w-5 text-red-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Total Cancellations</p>
              <p className="text-2xl font-bold text-gray-900">{totalCount}</p>
            </div>
          </div>
        </div>
        <div className="rounded-lg border border-gray-200 bg-white p-5">
          <div className="flex items-center gap-3">
            <div className="rounded-lg bg-amber-100 p-2.5">
              <DollarSign className="h-5 w-5 text-amber-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Value Lost</p>
              <p className="text-2xl font-bold text-gray-900">
                ${valueLost.toLocaleString("en-US", { minimumFractionDigits: 2 })}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className="overflow-hidden rounded-lg border border-gray-200 bg-white">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Order ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Customer</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Amount</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Reason</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.map((item) => (
              <tr key={item.id} className="hover:bg-gray-50">
                <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-blue-600">{item.orderNumber}</td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">{item.customerName || "Unknown"}</td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                  ${item.totalAmount?.toLocaleString("en-US", { minimumFractionDigits: 2 })}
                </td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">{item.cancellationReason || "N/A"}</td>
                <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                  {new Date(item.updatedAt || item.createdAt).toLocaleDateString()}
                </td>
                <td className="whitespace-nowrap px-6 py-4">
                  <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusBadge[item.status] || "bg-gray-100 text-gray-700"}`}>
                    {item.status}
                  </span>
                </td>
              </tr>
            ))}
            {orders.length === 0 && (
              <tr>
                <td colSpan={6} className="px-6 py-8 text-center text-gray-400">
                  No cancellations found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
