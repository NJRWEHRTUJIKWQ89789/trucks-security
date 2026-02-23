"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";

interface Activity { action: string; entityType: string; createdAt: string; }

export default function RecentActivity() {
  const [items, setItems] = useState<Activity[]>([]);

  useEffect(() => {
    gql<{ dashboardActivity: { items: Activity[] } }>(`{ dashboardActivity(page:1,perPage:8) { items { action entityType createdAt } } }`)
      .then((d) => setItems(d.dashboardActivity.items)).catch(() => {});
  }, []);

  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="font-semibold text-gray-900 mb-4">Recent Activity</h3>
      <div className="space-y-3">
        {items.length === 0 && <p className="text-sm text-gray-400">No recent activity</p>}
        {items.map((a, i) => (
          <div key={i} className="flex items-start gap-3 pb-3 border-b last:border-0">
            <div className="w-2 h-2 rounded-full bg-blue-500 mt-2 flex-shrink-0" />
            <div>
              <p className="text-sm text-gray-700">{a.action} <span className="text-gray-400">({a.entityType})</span></p>
              <p className="text-xs text-gray-400">{new Date(a.createdAt).toLocaleString()}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
