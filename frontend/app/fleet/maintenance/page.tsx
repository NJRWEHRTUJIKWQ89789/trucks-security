"use client";

import { useEffect, useState } from "react";
import { Wrench, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Maintenance { id: string; type: string; description: string; status: string; scheduledDate: string; completedDate: string; cost: number; mechanic: string; }

const statusColor: Record<string, string> = { scheduled: "bg-blue-100 text-blue-700", in_progress: "bg-amber-100 text-amber-700", completed: "bg-green-100 text-green-700", cancelled: "bg-red-100 text-red-700" };

export default function MaintenancePage() {
  const [records, setRecords] = useState<Maintenance[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ maintenanceRecords: { items: Maintenance[] } }>(`{ maintenanceRecords(page:1,perPage:50) { items { id type description status scheduledDate completedDate cost mechanic } } }`)
      .then((d) => setRecords(d.maintenanceRecords.items)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Maintenance Logs</h1>
      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Description</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Scheduled</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Cost</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Mechanic</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {records.map((r) => (
              <tr key={r.id} className="hover:bg-gray-50">
                <td className="px-6 py-4"><div className="flex items-center gap-2"><Wrench className="h-4 w-4 text-gray-400" /><span className="text-sm font-medium text-gray-900">{r.type}</span></div></td>
                <td className="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">{r.description}</td>
                <td className="px-6 py-4"><span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColor[r.status] || "bg-gray-100 text-gray-700"}`}>{r.status}</span></td>
                <td className="px-6 py-4 text-sm text-gray-500">{r.scheduledDate ? new Date(r.scheduledDate).toLocaleDateString() : "-"}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{r.cost ? `$${r.cost.toLocaleString()}` : "-"}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{r.mechanic}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
