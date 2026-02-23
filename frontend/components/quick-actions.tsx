import {
  Package,
  Search,
  Truck,
  UserPlus,
  Building,
  TrendingUp,
  Settings,
  HelpCircle,
} from "lucide-react";

const actions = [
  {
    title: "Create Shipment",
    description: "Add a new shipment to the system",
    icon: Package,
    color: "blue",
  },
  {
    title: "Track Package",
    description: "Track existing shipments",
    icon: Search,
    color: "green",
  },
  {
    title: "Add Vehicle",
    description: "Register a new vehicle",
    icon: Truck,
    color: "purple",
  },
  {
    title: "Add Vendor",
    description: "Register a new vendor",
    icon: UserPlus,
    color: "amber",
  },
  {
    title: "Inventory Check",
    description: "Check warehouse inventory",
    icon: Building,
    color: "indigo",
  },
  {
    title: "Generate Report",
    description: "Create performance reports",
    icon: TrendingUp,
    color: "orange",
  },
  {
    title: "System Settings",
    description: "Configure system preferences",
    icon: Settings,
    color: "gray",
  },
  {
    title: "Documentation",
    description: "Access help and guides",
    icon: HelpCircle,
    color: "teal",
  },
];

const colorClasses: Record<string, { bg: string; text: string }> = {
  blue: { bg: "bg-blue-50", text: "text-blue-500" },
  green: { bg: "bg-green-50", text: "text-green-500" },
  purple: { bg: "bg-purple-50", text: "text-purple-500" },
  amber: { bg: "bg-amber-50", text: "text-amber-500" },
  indigo: { bg: "bg-indigo-50", text: "text-indigo-500" },
  orange: { bg: "bg-orange-50", text: "text-orange-500" },
  gray: { bg: "bg-gray-50", text: "text-gray-500" },
  teal: { bg: "bg-teal-50", text: "text-teal-500" },
};

export default function QuickActions() {
  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <h3 className="text-lg font-semibold mb-4">Quick Actions</h3>
      <div className="grid grid-cols-2 gap-3">
        {actions.map((action) => {
          const Icon = action.icon;
          const colors = colorClasses[action.color];
          return (
            <div
              key={action.title}
              className="flex items-center gap-3 p-3 rounded-lg border hover:bg-gray-50 cursor-pointer transition-colors"
            >
              <div
                className={`w-10 h-10 rounded-full ${colors.bg} flex items-center justify-center ${colors.text}`}
              >
                <Icon className="w-5 h-5" />
              </div>
              <div>
                <div className="text-sm font-medium">{action.title}</div>
                <div className="text-xs text-gray-500">
                  {action.description}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
