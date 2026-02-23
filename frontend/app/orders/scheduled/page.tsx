"use client";

import { useEffect, useState } from "react";
import { Calendar, Truck, MapPin, User, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Order {
  id: string;
  orderNumber: string;
  customerName: string;
  status: string;
  totalAmount: number;
  shipmentId: string;
  scheduledDate: string;
  createdAt: string;
}

const statusBadge: Record<string, string> = {
  scheduled: "bg-green-100 text-green-700",
  pending: "bg-amber-100 text-amber-700",
  confirmed: "bg-green-100 text-green-700",
  rescheduled: "bg-blue-100 text-blue-700",
};

function groupByDate(orders: Order[]): { date: string; orders: Order[] }[] {
  const groups: Record<string, Order[]> = {};
  for (const order of orders) {
    const dateKey = order.scheduledDate
      ? new Date(order.scheduledDate).toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" })
      : new Date(order.createdAt).toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" });
    if (!groups[dateKey]) groups[dateKey] = [];
    groups[dateKey].push(order);
  }
  return Object.entries(groups)
    .sort(([a], [b]) => new Date(a).getTime() - new Date(b).getTime())
    .map(([date, orders]) => ({ date, orders }));
}

export default function ScheduledDeliveriesPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ scheduledOrders: { items: Order[] } }>(`{
      scheduledOrders(page: 1, perPage: 50) {
        items {
          id orderNumber customerName status totalAmount shipmentId scheduledDate createdAt
        }
      }
    }`)
      .then((d) => setOrders(d.scheduledOrders.items))
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

  const grouped = groupByDate(orders);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Scheduled Deliveries</h1>

      {orders.length === 0 && (
        <div className="rounded-lg border border-gray-200 bg-white px-6 py-12 text-center text-gray-400">
          No scheduled deliveries found
        </div>
      )}

      <div className="space-y-8">
        {grouped.map((day) => (
          <section key={day.date}>
            <div className="mb-4 flex items-center gap-2">
              <Calendar className="h-5 w-5 text-gray-400" />
              <h2 className="text-lg font-semibold text-gray-900">{day.date}</h2>
            </div>

            <div className="space-y-3">
              {day.orders.map((order) => (
                <div
                  key={order.id}
                  className="flex items-center justify-between rounded-lg border border-gray-200 bg-white px-5 py-4 hover:shadow-sm"
                >
                  <div className="flex items-center gap-6">
                    <span className="w-20 text-sm font-medium text-gray-900">
                      {order.scheduledDate
                        ? new Date(order.scheduledDate).toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" })
                        : "--:--"}
                    </span>
                    <div className="flex items-center gap-2">
                      <Truck className="h-4 w-4 text-gray-400" />
                      <span className="text-sm font-medium text-blue-600">{order.orderNumber}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <MapPin className="h-4 w-4 text-gray-400" />
                      <span className="text-sm text-gray-700">
                        {order.shipmentId ? `Shipment ${order.shipmentId.slice(0, 8)}` : "No shipment"}
                      </span>
                    </div>
                    <div className="flex items-center gap-2">
                      <User className="h-4 w-4 text-gray-400" />
                      <span className="text-sm text-gray-500">{order.customerName || "Unknown"}</span>
                    </div>
                  </div>
                  <span
                    className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusBadge[order.status] || "bg-gray-100 text-gray-700"}`}
                  >
                    {order.status}
                  </span>
                </div>
              ))}
            </div>
          </section>
        ))}
      </div>
    </div>
  );
}
