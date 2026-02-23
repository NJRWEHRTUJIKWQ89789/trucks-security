import { ScrollText } from "lucide-react";

const logs = [
  { timestamp: "2026-02-22 14:32:05", user: "admin@cargomax.io", action: "Updated", resource: "Settings / General", ip: "192.168.1.10" },
  { timestamp: "2026-02-22 14:18:41", user: "m.chen@cargomax.io", action: "Created", resource: "Shipment #SHP-4521", ip: "10.0.0.45" },
  { timestamp: "2026-02-22 13:55:12", user: "j.smith@cargomax.io", action: "Deleted", resource: "Driver / R. Patel", ip: "172.16.0.22" },
  { timestamp: "2026-02-22 12:40:33", user: "admin@cargomax.io", action: "Modified", resource: "Role / Dispatcher", ip: "192.168.1.10" },
  { timestamp: "2026-02-22 11:22:17", user: "s.jones@cargomax.io", action: "Exported", resource: "Report / February Fleet", ip: "10.0.0.88" },
  { timestamp: "2026-02-22 10:15:44", user: "m.chen@cargomax.io", action: "Assigned", resource: "Vehicle #VH-312 → Route 7A", ip: "10.0.0.45" },
  { timestamp: "2026-02-21 17:48:29", user: "admin@cargomax.io", action: "Enabled", resource: "2FA / Organization", ip: "192.168.1.10" },
  { timestamp: "2026-02-21 16:30:11", user: "j.smith@cargomax.io", action: "Created", resource: "Warehouse / Denver Hub", ip: "172.16.0.22" },
  { timestamp: "2026-02-21 14:05:52", user: "s.jones@cargomax.io", action: "Uploaded", resource: "Document / Insurance-2026.pdf", ip: "10.0.0.88" },
  { timestamp: "2026-02-21 09:12:08", user: "system", action: "Scheduled", resource: "Maintenance / VH-208", ip: "127.0.0.1" },
];

const actionColors: Record<string, string> = {
  Created: "text-green-600",
  Updated: "text-blue-600",
  Modified: "text-blue-600",
  Deleted: "text-red-600",
  Exported: "text-purple-600",
  Assigned: "text-orange-600",
  Enabled: "text-emerald-600",
  Uploaded: "text-cyan-600",
  Scheduled: "text-gray-600",
};

export default function LogsPage() {
  return (
    <div className="">
      <div className="flex items-center gap-3 mb-6">
        <ScrollText className="w-6 h-6 text-blue-600" />
        <h1 className="text-2xl font-bold">Audit Logs</h1>
      </div>

      <div className="bg-white rounded-lg border overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50 border-b text-sm text-gray-600">
              <th className="text-left px-4 py-3 font-semibold">Timestamp</th>
              <th className="text-left px-4 py-3 font-semibold">User</th>
              <th className="text-left px-4 py-3 font-semibold">Action</th>
              <th className="text-left px-4 py-3 font-semibold">Resource</th>
              <th className="text-left px-4 py-3 font-semibold">IP Address</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log, i) => (
              <tr key={i} className="border-b last:border-0 hover:bg-gray-50">
                <td className="px-4 py-3 text-sm font-mono text-gray-500">{log.timestamp}</td>
                <td className="px-4 py-3 text-sm">{log.user}</td>
                <td className={`px-4 py-3 text-sm font-medium ${actionColors[log.action] ?? "text-gray-700"}`}>
                  {log.action}
                </td>
                <td className="px-4 py-3 text-sm">{log.resource}</td>
                <td className="px-4 py-3 text-sm font-mono text-gray-400">{log.ip}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
