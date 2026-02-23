"use client";

import { useEffect, useState } from "react";
import { Users, Loader2, Mail, Phone, Star, DollarSign, Package } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Client { id: string; companyName: string; contactPerson: string; email: string; phone: string; industry: string; totalShipments: number; totalSpent: number; satisfactionRating: number; status: string; }

export default function ClientsPage() {
  const [clients, setClients] = useState<Client[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ clients: { items: Client[] } }>(`{ clients(page:1,perPage:50) { items { id companyName contactPerson email phone industry totalShipments totalSpent satisfactionRating status } } }`)
      .then((d) => setClients(d.clients.items)).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="flex justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-blue-600" /></div>;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Clients List</h1>
      <div className="grid grid-cols-2 gap-4">
        {clients.map((c) => (
          <div key={c.id} className="rounded-xl border bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-3"><div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center"><Users className="h-5 w-5 text-blue-600" /></div><div><p className="font-semibold text-gray-900">{c.companyName}</p><p className="text-xs text-gray-500">{c.industry} &middot; {c.contactPerson}</p></div></div>
              <span className={`text-xs px-2 py-0.5 rounded-full ${c.status === "active" ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-700"}`}>{c.status}</span>
            </div>
            <div className="space-y-1 text-sm text-gray-600">
              <div className="flex items-center gap-2"><Mail className="h-3.5 w-3.5" />{c.email}</div>
              <div className="flex items-center gap-2"><Phone className="h-3.5 w-3.5" />{c.phone}</div>
            </div>
            <div className="flex items-center gap-4 mt-3 pt-3 border-t text-sm">
              <span className="flex items-center gap-1 text-gray-500"><Package className="h-3.5 w-3.5" />{c.totalShipments} shipments</span>
              <span className="flex items-center gap-1 text-gray-500"><DollarSign className="h-3.5 w-3.5" />${c.totalSpent?.toLocaleString()}</span>
              <span className="flex items-center gap-1 text-amber-500 ml-auto"><Star className="h-3.5 w-3.5 fill-current" />{c.satisfactionRating?.toFixed(1)}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
