import { Ticket, Plus } from "lucide-react";

const tickets = [
  { id: "TKT-1041", subject: "Unable to generate monthly report", priority: "High", status: "Open", created: "Feb 20, 2026", updated: "Feb 22, 2026" },
  { id: "TKT-1040", subject: "Driver app GPS not updating", priority: "Critical", status: "In Progress", created: "Feb 19, 2026", updated: "Feb 22, 2026" },
  { id: "TKT-1038", subject: "Invoice amount mismatch on order #8832", priority: "Medium", status: "Open", created: "Feb 18, 2026", updated: "Feb 20, 2026" },
  { id: "TKT-1035", subject: "Cannot assign driver to new route", priority: "High", status: "In Progress", created: "Feb 16, 2026", updated: "Feb 19, 2026" },
  { id: "TKT-1032", subject: "API rate limit too restrictive", priority: "Low", status: "Resolved", created: "Feb 14, 2026", updated: "Feb 17, 2026" },
  { id: "TKT-1029", subject: "Warehouse map not loading on Safari", priority: "Medium", status: "Resolved", created: "Feb 12, 2026", updated: "Feb 15, 2026" },
];

const priorityStyles: Record<string, string> = {
  Critical: "bg-red-100 text-red-700",
  High: "bg-orange-100 text-orange-700",
  Medium: "bg-yellow-100 text-yellow-700",
  Low: "bg-gray-100 text-gray-600",
};

const statusStyles: Record<string, string> = {
  Open: "bg-blue-100 text-blue-700",
  "In Progress": "bg-amber-100 text-amber-700",
  Resolved: "bg-green-100 text-green-700",
};

export default function TicketsPage() {
  return (
    <div className="">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Ticket className="w-6 h-6 text-blue-600" />
          <h1 className="text-2xl font-bold">Support Tickets</h1>
        </div>
        <button className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors text-sm">
          <Plus className="w-4 h-4" />
          New Ticket
        </button>
      </div>

      <div className="bg-white rounded-lg border overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50 border-b text-sm text-gray-600">
              <th className="text-left px-4 py-3 font-semibold">Ticket #</th>
              <th className="text-left px-4 py-3 font-semibold">Subject</th>
              <th className="text-left px-4 py-3 font-semibold">Priority</th>
              <th className="text-left px-4 py-3 font-semibold">Status</th>
              <th className="text-left px-4 py-3 font-semibold">Created</th>
              <th className="text-left px-4 py-3 font-semibold">Updated</th>
            </tr>
          </thead>
          <tbody>
            {tickets.map((t) => (
              <tr key={t.id} className="border-b last:border-0 hover:bg-gray-50 cursor-pointer">
                <td className="px-4 py-3 text-sm font-mono text-blue-600">{t.id}</td>
                <td className="px-4 py-3 text-sm font-medium">{t.subject}</td>
                <td className="px-4 py-3">
                  <span className={`text-xs font-medium px-2 py-1 rounded-full ${priorityStyles[t.priority]}`}>
                    {t.priority}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <span className={`text-xs font-medium px-2 py-1 rounded-full ${statusStyles[t.status]}`}>
                    {t.status}
                  </span>
                </td>
                <td className="px-4 py-3 text-sm text-gray-500">{t.created}</td>
                <td className="px-4 py-3 text-sm text-gray-500">{t.updated}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
