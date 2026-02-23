"use client";

import { useEffect, useState } from "react";
import { Store, Loader2, Star, Mail, Phone } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Vendor { id: string; name: string; contactPerson: string; email: string; phone: string; category: string; rating: number; status: string; }

const statusColor: Record<string, string> = { active: "bg-green-100 text-green-700", inactive: "bg-gray-100 text-gray-700", suspended: "bg-red-100 text-red-700" };

export default function VendorsPage() {
  const [vendors, setVendors] = useState<Vendor[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ vendors: { items: Vendor[] } }>(`{ vendors(page:1,perPage:50) { items { id name contactPerson email phone category rating status } } }`)
      .then((d) => setVendors(d.vendors.items)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Vendor Directory</h1>
      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Vendor</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Contact</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Category</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Rating</th>
              <th className="px-6 py-3 text-left text-xs font-medium uppercase text-gray-500">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {vendors.map((v) => (
              <tr key={v.id} className="hover:bg-gray-50">
                <td className="px-6 py-4"><div className="flex items-center gap-3"><Store className="h-5 w-5 text-gray-400" /><div><p className="text-sm font-medium text-gray-900">{v.name}</p><p className="text-xs text-gray-500">{v.contactPerson}</p></div></div></td>
                <td className="px-6 py-4 text-sm"><div className="flex items-center gap-1 text-gray-600"><Mail className="h-3.5 w-3.5" />{v.email}</div><div className="flex items-center gap-1 text-gray-500 text-xs mt-0.5"><Phone className="h-3 w-3" />{v.phone}</div></td>
                <td className="px-6 py-4 text-sm text-gray-700">{v.category}</td>
                <td className="px-6 py-4"><span className="flex items-center gap-1 text-amber-500"><Star className="h-4 w-4 fill-current" />{v.rating?.toFixed(1)}</span></td>
                <td className="px-6 py-4"><span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColor[v.status] || "bg-gray-100 text-gray-700"}`}>{v.status}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
