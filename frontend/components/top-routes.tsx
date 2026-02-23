"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";

interface Shipment { origin: string; destination: string; }

export default function TopRoutes() {
  const [routes, setRoutes] = useState<{ name: string; count: string }[]>([]);

  useEffect(() => {
    gql<{ shipments: { items: Shipment[] } }>(`{ shipments(page:1,perPage:100) { items { origin destination } } }`)
      .then((d) => {
        const counts: Record<string, number> = {};
        for (const s of d.shipments.items) {
          const key = `${s.origin || "?"} → ${s.destination || "?"}`;
          counts[key] = (counts[key] || 0) + 1;
        }
        const sorted = Object.entries(counts).sort((a, b) => b[1] - a[1]).slice(0, 5);
        setRoutes(sorted.map(([name, c]) => ({ name, count: `${c} shipments` })));
      }).catch(() => {});
  }, []);

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="font-semibold text-gray-900 mb-4">Top Routes</h3>
      <div className="space-y-3">
        {routes.map((r, i) => (
          <div key={i} className="flex items-center justify-between">
            <span className="text-sm text-gray-700">{r.name}</span>
            <span className="text-xs text-gray-400">{r.count}</span>
          </div>
        ))}
        {routes.length === 0 && <p className="text-sm text-gray-400">Loading...</p>}
      </div>
    </div>
  );
}
