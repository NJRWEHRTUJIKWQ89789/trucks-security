"use client";

import { useEffect, useState } from "react";
import { Package, AlertTriangle, TrendingUp, CheckCircle, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Warehouse {
  id: string;
  name: string;
}

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

interface InventoryRow extends InventoryItem {
  warehouseName: string;
}

function deriveDisplayStatus(item: InventoryItem): string {
  if (item.quantity <= 0) return "Critical";
  if (item.quantity <= item.minQuantity * 0.5) return "Critical";
  if (item.quantity <= item.minQuantity) return "Low Stock";
  if (item.quantity > item.minQuantity * 3) return "Overstocked";
  return "Optimal";
}

function statusBadge(status: string) {
  switch (status) {
    case "Optimal":
      return <span className="inline-flex items-center rounded-full bg-green-50 px-2.5 py-0.5 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20">Optimal</span>;
    case "Low Stock":
      return <span className="inline-flex items-center rounded-full bg-red-50 px-2.5 py-0.5 text-xs font-medium text-red-700 ring-1 ring-inset ring-red-600/20">Low Stock</span>;
    case "Overstocked":
      return <span className="inline-flex items-center rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-700 ring-1 ring-inset ring-amber-600/20">Overstocked</span>;
    case "Critical":
      return <span className="inline-flex items-center rounded-full bg-red-50 px-2.5 py-0.5 text-xs font-bold text-red-700 ring-1 ring-inset ring-red-600/30">Critical</span>;
    default:
      return null;
  }
}

export default function InventoryLevelsPage() {
  const [inventory, setInventory] = useState<InventoryRow[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchInventory() {
      try {
        const warehouseData = await gql<{ warehouses: { items: Warehouse[] } }>(
          `{ warehouses(page:1, perPage:100) { items { id name } } }`
        );
        const warehouses = warehouseData.warehouses.items;
        const warehouseMap = new Map(warehouses.map((w) => [w.id, w.name]));

        const allItems: InventoryRow[] = [];
        await Promise.all(
          warehouses.map(async (w) => {
            const data = await gql<{ inventoryItems: { items: InventoryItem[] } }>(
              `query($wid: String!) { inventoryItems(warehouseId: $wid, page: 1, perPage: 200) { items { id sku name category warehouseId quantity minQuantity status } } }`,
              { wid: w.id }
            );
            for (const item of data.inventoryItems.items) {
              allItems.push({
                ...item,
                warehouseName: warehouseMap.get(item.warehouseId) || "Unknown",
              });
            }
          })
        );

        setInventory(allItems);
      } catch {
        // silently handle errors
      } finally {
        setLoading(false);
      }
    }
    fetchInventory();
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  const totalItems = inventory.reduce((sum, item) => sum + item.quantity, 0);
  const lowStockCount = inventory.filter((i) => deriveDisplayStatus(i) === "Low Stock").length;
  const criticalCount = inventory.filter((i) => deriveDisplayStatus(i) === "Critical").length;
  const overstockedCount = inventory.filter((i) => deriveDisplayStatus(i) === "Overstocked").length;
  const optimalCount = inventory.filter((i) => deriveDisplayStatus(i) === "Optimal").length;

  const summaryCards = [
    { label: "Total Items", value: totalItems.toLocaleString(), icon: Package, color: "bg-blue-50 text-blue-600", border: "border-blue-200" },
    { label: "Low Stock", value: String(lowStockCount + criticalCount), icon: AlertTriangle, color: "bg-red-50 text-red-600", border: "border-red-200" },
    { label: "Overstocked", value: String(overstockedCount), icon: TrendingUp, color: "bg-amber-50 text-amber-600", border: "border-amber-200" },
    { label: "Optimal", value: String(optimalCount), icon: CheckCircle, color: "bg-green-50 text-green-600", border: "border-green-200" },
  ];

  return (
    <div className="mx-auto max-w-7xl">
        <h1 className="mb-8 text-3xl font-bold text-gray-900">Inventory Levels</h1>

        <div className="mb-8 grid grid-cols-4 gap-4">
          {summaryCards.map((card) => (
            <div key={card.label} className={`rounded-xl border bg-white p-6 shadow-sm ${card.border}`}>
              <div className="flex items-center gap-3">
                <div className={`rounded-lg p-2.5 ${card.color}`}>
                  <card.icon className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">{card.label}</p>
                  <p className="text-2xl font-bold text-gray-900">{card.value}</p>
                </div>
              </div>
            </div>
          ))}
        </div>

        {inventory.length === 0 ? (
          <div className="rounded-xl border border-gray-200 bg-white p-12 text-center shadow-sm">
            <Package className="mx-auto h-12 w-12 text-gray-300" />
            <p className="mt-4 text-gray-500">No inventory items found.</p>
          </div>
        ) : (
          <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50">
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">SKU</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Name</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Category</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Warehouse</th>
                  <th className="px-6 py-3.5 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">Qty</th>
                  <th className="px-6 py-3.5 text-right text-xs font-semibold uppercase tracking-wider text-gray-500">Min Level</th>
                  <th className="px-6 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {inventory.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-50">
                    <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">{item.sku}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-700">{item.name}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">{item.category}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">{item.warehouseName}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium text-gray-900">{item.quantity.toLocaleString()}</td>
                    <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-gray-500">{item.minQuantity.toLocaleString()}</td>
                    <td className="whitespace-nowrap px-6 py-4">{statusBadge(deriveDisplayStatus(item))}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
  );
}
