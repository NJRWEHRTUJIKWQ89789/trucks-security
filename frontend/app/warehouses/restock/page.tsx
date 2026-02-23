"use client";

import { useEffect, useState } from "react";
import { Plus, Loader2, PackageCheck } from "lucide-react";
import { gql } from "@/lib/graphql";

interface InventoryItem {
  id: string;
  sku: string;
  name: string;
  category: string;
  warehouseId: string;
  quantity: number;
  minQuantity: number;
  status: string;
}

interface Warehouse {
  id: string;
  name: string;
}

interface LowStockRow extends InventoryItem {
  warehouseName: string;
  deficit: number;
}

function priorityFromDeficit(item: LowStockRow): string {
  const ratio = item.minQuantity > 0 ? item.quantity / item.minQuantity : 0;
  if (ratio <= 0) return "Urgent";
  if (ratio <= 0.25) return "Urgent";
  if (ratio <= 0.5) return "High";
  if (ratio <= 0.75) return "Normal";
  return "Low";
}

function priorityBadge(priority: string) {
  const styles: Record<string, string> = {
    Urgent: "bg-red-50 text-red-700 ring-red-600/20",
    High: "bg-amber-50 text-amber-700 ring-amber-600/20",
    Normal: "bg-blue-50 text-blue-700 ring-blue-600/20",
    Low: "bg-gray-50 text-gray-600 ring-gray-500/20",
  };
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset ${styles[priority]}`}>
      {priority}
    </span>
  );
}

function statusBadge(status: string) {
  const styles: Record<string, string> = {
    Pending: "bg-amber-50 text-amber-700 ring-amber-600/20",
    Restocked: "bg-green-50 text-green-700 ring-green-600/20",
  };
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset ${styles[status] || "bg-gray-50 text-gray-600 ring-gray-500/20"}`}>
      {status}
    </span>
  );
}

export default function RestockRequestsPage() {
  const [items, setItems] = useState<LowStockRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [restocking, setRestocking] = useState<string | null>(null);

  async function fetchData() {
    try {
      const [lowStockData, warehouseData] = await Promise.all([
        gql<{ lowStockItems: InventoryItem[] }>(
          `{ lowStockItems { id sku name category warehouseId quantity minQuantity status } }`
        ),
        gql<{ warehouses: { items: Warehouse[] } }>(
          `{ warehouses(page:1, perPage:100) { items { id name } } }`
        ),
      ]);

      const warehouseMap = new Map(warehouseData.warehouses.items.map((w) => [w.id, w.name]));

      const rows: LowStockRow[] = (lowStockData.lowStockItems || []).map((item) => ({
        ...item,
        warehouseName: warehouseMap.get(item.warehouseId) || "Unknown",
        deficit: Math.max(0, item.minQuantity - item.quantity),
      }));

      rows.sort((a, b) => {
        const order = { Urgent: 0, High: 1, Normal: 2, Low: 3 };
        return (order[priorityFromDeficit(a) as keyof typeof order] ?? 4) - (order[priorityFromDeficit(b) as keyof typeof order] ?? 4);
      });

      setItems(rows);
    } catch {
      // silently handle errors
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchData();
  }, []);

  async function handleRestock(item: LowStockRow) {
    if (restocking) return;
    setRestocking(item.id);
    try {
      await gql<{ restockItem: InventoryItem }>(
        `mutation($id: String!, $quantity: Int!) { restockItem(id: $id, quantity: $quantity) { id quantity status } }`,
        { id: item.id, quantity: item.deficit }
      );
      await fetchData();
    } catch {
      // silently handle errors
    } finally {
      setRestocking(null);
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-7xl">
        <div className="mb-8 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Restock Requests</h1>
          <button
            onClick={() => fetchData()}
            className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-blue-700"
          >
            <Plus className="h-4 w-4" />
            Refresh
          </button>
        </div>

        {items.length === 0 ? (
          <div className="rounded-xl border border-gray-200 bg-white p-12 text-center shadow-sm">
            <PackageCheck className="mx-auto h-12 w-12 text-green-300" />
            <p className="mt-4 text-gray-500">All inventory levels are healthy. No restock needed.</p>
          </div>
        ) : (
          <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50">
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">SKU</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Item</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Warehouse</th>
                  <th className="px-6 py-3.5 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">Current Qty</th>
                  <th className="px-6 py-3.5 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">Min Level</th>
                  <th className="px-6 py-3.5 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">Deficit</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Priority</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Status</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {items.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-50">
                    <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">{item.sku}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-700">{item.name}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">{item.warehouseName}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium text-gray-900">{item.quantity.toLocaleString()}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-gray-500">{item.minQuantity.toLocaleString()}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium text-red-600">{item.deficit.toLocaleString()}</td>
                    <td className="whitespace-nowrap px-6 py-4">{priorityBadge(priorityFromDeficit(item))}</td>
                    <td className="whitespace-nowrap px-6 py-4">{statusBadge(item.status === "in_stock" ? "Restocked" : "Pending")}</td>
                    <td className="whitespace-nowrap px-6 py-4">
                      <button
                        onClick={() => handleRestock(item)}
                        disabled={restocking === item.id || item.deficit === 0}
                        className="inline-flex items-center rounded-md bg-blue-600 px-3 py-1.5 text-xs font-medium text-white shadow-sm hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {restocking === item.id ? (
                          <Loader2 className="h-3 w-3 animate-spin" />
                        ) : (
                          `Restock +${item.deficit.toLocaleString()}`
                        )}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
  );
}
