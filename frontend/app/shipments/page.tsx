"use client";

import { useState, useEffect } from "react";
import { Search, Plus, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Shipment {
  id: string;
  trackingNumber: string;
  origin: string;
  destination: string;
  status: string;
  carrier: string;
  weight: number;
  customerName: string;
  createdAt: string;
}

const statusColors: Record<string, string> = {
  in_transit: "bg-blue-100 text-blue-700",
  delivered: "bg-green-100 text-green-700",
  delayed: "bg-red-100 text-red-700",
  pending: "bg-amber-100 text-amber-700",
};

export default function ShipmentsPage() {
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("All");

  useEffect(() => {
    const status = statusFilter === "All" ? undefined : statusFilter.toLowerCase().replace(" ", "_");
    gql<{ shipments: { items: Shipment[] } }>(`query($page:Int,$perPage:Int,$status:String){shipments(page:$page,perPage:$perPage,status:$status){items{id trackingNumber origin destination status carrier weight customerName createdAt}}}`, { page: 1, perPage: 50, status })
      .then((d) => setShipments(d.shipments.items))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [statusFilter]);

  const filtered = shipments.filter((s) => {
    const q = search.toLowerCase();
    return !q || s.trackingNumber?.toLowerCase().includes(q) || s.origin?.toLowerCase().includes(q) || s.destination?.toLowerCase().includes(q) || s.customerName?.toLowerCase().includes(q);
  });

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold text-gray-900">All Shipments</h1>
        <button className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors">
          <Plus className="w-4 h-4" /> New Shipment
        </button>
      </div>
      <div className="flex items-center gap-4 mb-6">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input type="text" placeholder="Search shipments..." value={search} onChange={(e) => setSearch(e.target.value)} className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
        </div>
        <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="border border-gray-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500">
          <option>All</option>
          <option>Pending</option>
          <option>In Transit</option>
          <option>Delivered</option>
          <option>Delayed</option>
        </select>
      </div>
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50 border-b border-gray-200">
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Tracking #</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Origin</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Destination</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Carrier</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Weight</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Customer</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {filtered.map((s) => (
              <tr key={s.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-6 py-4 text-sm font-medium text-blue-600">{s.trackingNumber}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{s.origin}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{s.destination}</td>
                <td className="px-6 py-4"><span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium ${statusColors[s.status] || "bg-gray-100 text-gray-700"}`}>{s.status}</span></td>
                <td className="px-6 py-4 text-sm text-gray-700">{s.carrier}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{s.weight ? `${s.weight} kg` : "-"}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{s.customerName}</td>
              </tr>
            ))}
            {filtered.length === 0 && <tr><td colSpan={7} className="px-6 py-8 text-center text-gray-400">No shipments found</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  );
}
