import {
  TrendingUp, TrendingDown, Package, Truck, Users, DollarSign,
  AlertCircle, AlertTriangle, CheckCircle, Info, ChevronRight,
} from "lucide-react";

const stats = [
  { label: "Total Shipments", value: "2,847", change: "+12.5%", up: true, icon: Package, color: "bg-blue-50 text-blue-600" },
  { label: "Active Vehicles", value: "156", change: "+3.2%", up: true, icon: Truck, color: "bg-green-50 text-green-600" },
  { label: "Total Drivers", value: "89", change: "-2.1%", up: false, icon: Users, color: "bg-orange-50 text-orange-600" },
  { label: "Revenue", value: "$1.2M", change: "+8.7%", up: true, icon: DollarSign, color: "bg-purple-50 text-purple-600" },
];

const progressBars = [
  { label: "Fleet Utilization", value: 78, color: "bg-blue-500" },
  { label: "On-time Delivery", value: 94, color: "bg-green-500" },
  { label: "Fuel Efficiency", value: 65, color: "bg-orange-500" },
  { label: "Driver Satisfaction", value: 88, color: "bg-purple-500" },
];

const tableData = [
  { id: "SHP-4521", origin: "San Francisco", dest: "Denver", status: "In Transit", eta: "Feb 23" },
  { id: "SHP-4520", origin: "Los Angeles", dest: "Phoenix", status: "Delivered", eta: "Feb 22" },
  { id: "SHP-4519", origin: "Seattle", dest: "Portland", status: "Delayed", eta: "Feb 24" },
  { id: "SHP-4518", origin: "Chicago", dest: "Detroit", status: "In Transit", eta: "Feb 23" },
];

const badges = [
  "Active", "Pending", "Cancelled", "In Transit", "Delivered", "Delayed",
  "Urgent", "Scheduled", "Draft", "Completed", "On Hold", "Processing",
];

const badgeColors: Record<string, string> = {
  Active: "bg-green-100 text-green-700", Pending: "bg-yellow-100 text-yellow-700",
  Cancelled: "bg-red-100 text-red-700", "In Transit": "bg-blue-100 text-blue-700",
  Delivered: "bg-emerald-100 text-emerald-700", Delayed: "bg-orange-100 text-orange-700",
  Urgent: "bg-red-500 text-white", Scheduled: "bg-indigo-100 text-indigo-700",
  Draft: "bg-gray-100 text-gray-600", Completed: "bg-teal-100 text-teal-700",
  "On Hold": "bg-amber-100 text-amber-700", Processing: "bg-cyan-100 text-cyan-700",
};

const alerts = [
  { type: "info", icon: Info, title: "System update scheduled", desc: "Maintenance window: Feb 25, 2:00-4:00 AM PST", border: "border-blue-200 bg-blue-50", iconColor: "text-blue-500" },
  { type: "success", icon: CheckCircle, title: "Backup completed", desc: "All data successfully backed up at 6:00 AM.", border: "border-green-200 bg-green-50", iconColor: "text-green-500" },
  { type: "warning", icon: AlertTriangle, title: "Storage nearing capacity", desc: "85% of allocated storage has been used.", border: "border-yellow-200 bg-yellow-50", iconColor: "text-yellow-500" },
  { type: "error", icon: AlertCircle, title: "API rate limit exceeded", desc: "Reduce request frequency or upgrade your plan.", border: "border-red-200 bg-red-50", iconColor: "text-red-500" },
];

export default function WidgetsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Widget Showcase</h1>

      {/* Stat Cards */}
      <section>
        <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Stat Cards</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {stats.map((s) => (
            <div key={s.label} className="bg-white rounded-lg border p-4">
              <div className="flex items-center justify-between mb-3">
                <div className={`w-9 h-9 rounded-lg ${s.color} flex items-center justify-center`}>
                  <s.icon className="w-5 h-5" />
                </div>
                <span className={`flex items-center gap-1 text-xs font-medium ${s.up ? "text-green-600" : "text-red-500"}`}>
                  {s.up ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
                  {s.change}
                </span>
              </div>
              <p className="text-2xl font-bold">{s.value}</p>
              <p className="text-xs text-gray-500 mt-1">{s.label}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Progress Bars */}
      <section>
        <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Progress Bars</h2>
        <div className="bg-white rounded-lg border p-4 space-y-4">
          {progressBars.map((p) => (
            <div key={p.label}>
              <div className="flex justify-between text-sm mb-1">
                <span className="font-medium">{p.label}</span>
                <span className="text-gray-500">{p.value}%</span>
              </div>
              <div className="w-full h-2 bg-gray-100 rounded-full">
                <div className={`h-2 rounded-full ${p.color}`} style={{ width: `${p.value}%` }} />
              </div>
            </div>
          ))}
        </div>
      </section>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Mini Chart */}
        <section>
          <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Mini Chart</h2>
          <div className="bg-white rounded-lg border p-4">
            <div className="flex items-end justify-between h-32 gap-2 px-2">
              {[40, 65, 45, 80, 55, 90, 70, 85, 60, 95, 75, 88].map((h, i) => (
                <div key={i} className="flex-1 flex flex-col items-center gap-1">
                  <div className="w-full bg-blue-500 rounded-t opacity-80 hover:opacity-100 transition-opacity" style={{ height: `${h}%` }} />
                  <span className="text-[9px] text-gray-400">{["J","F","M","A","M","J","J","A","S","O","N","D"][i]}</span>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Data Table */}
        <section>
          <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Data Table</h2>
          <div className="bg-white rounded-lg border overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-gray-50 border-b">
                  <th className="text-left px-3 py-2 font-semibold text-gray-600">ID</th>
                  <th className="text-left px-3 py-2 font-semibold text-gray-600">Route</th>
                  <th className="text-left px-3 py-2 font-semibold text-gray-600">Status</th>
                  <th className="text-left px-3 py-2 font-semibold text-gray-600">ETA</th>
                </tr>
              </thead>
              <tbody>
                {tableData.map((r) => (
                  <tr key={r.id} className="border-b last:border-0 hover:bg-gray-50">
                    <td className="px-3 py-2 font-mono text-blue-600">{r.id}</td>
                    <td className="px-3 py-2">{r.origin} → {r.dest}</td>
                    <td className="px-3 py-2">
                      <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${badgeColors[r.status] ?? "bg-gray-100 text-gray-600"}`}>
                        {r.status}
                      </span>
                    </td>
                    <td className="px-3 py-2 text-gray-500">{r.eta}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      </div>

      {/* Badge Collection */}
      <section>
        <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Badge Collection</h2>
        <div className="bg-white rounded-lg border p-4 flex flex-wrap gap-2">
          {badges.map((b) => (
            <span key={b} className={`text-xs font-medium px-3 py-1 rounded-full ${badgeColors[b]}`}>{b}</span>
          ))}
        </div>
      </section>

      {/* Alert Variants */}
      <section>
        <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Alert Variants</h2>
        <div className="space-y-3">
          {alerts.map((a) => (
            <div key={a.type} className={`flex items-start gap-3 px-4 py-3 rounded-lg border ${a.border}`}>
              <a.icon className={`w-5 h-5 mt-0.5 shrink-0 ${a.iconColor}`} />
              <div>
                <p className="text-sm font-medium">{a.title}</p>
                <p className="text-xs text-gray-600 mt-0.5">{a.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Button Styles */}
      <section>
        <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Button Styles</h2>
        <div className="bg-white rounded-lg border p-4 flex flex-wrap gap-3">
          <button className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700">Primary</button>
          <button className="px-4 py-2 bg-gray-800 text-white rounded-lg text-sm hover:bg-gray-900">Secondary</button>
          <button className="px-4 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50">Outline</button>
          <button className="px-4 py-2 bg-green-600 text-white rounded-lg text-sm hover:bg-green-700">Success</button>
          <button className="px-4 py-2 bg-red-600 text-white rounded-lg text-sm hover:bg-red-700">Danger</button>
          <button className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg text-sm hover:bg-gray-200">Ghost</button>
          <button className="flex items-center gap-1 px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700">
            Next <ChevronRight className="w-4 h-4" />
          </button>
          <button className="px-4 py-2 bg-blue-100 text-blue-700 rounded-lg text-sm hover:bg-blue-200">Soft</button>
        </div>
      </section>

      {/* Form Input Examples */}
      <section>
        <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Form Inputs</h2>
        <div className="bg-white rounded-lg border p-4 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Text Input</label>
            <input type="text" placeholder="Enter text..." className="w-full border rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Select</label>
            <select className="w-full border rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none">
              <option>Option 1</option>
              <option>Option 2</option>
              <option>Option 3</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Disabled</label>
            <input type="text" disabled value="Read only" className="w-full border rounded-lg px-3 py-2 text-sm bg-gray-50 text-gray-400" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">With Error</label>
            <input type="text" defaultValue="Invalid value" className="w-full border border-red-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-red-500 focus:outline-none" />
            <p className="text-xs text-red-500 mt-1">This field is required.</p>
          </div>
        </div>
      </section>
    </div>
  );
}
