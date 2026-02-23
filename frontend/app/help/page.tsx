import { Search, Rocket, Package, Truck, CreditCard, Code, HelpCircle } from "lucide-react";

const cards = [
  { title: "Getting Started", description: "Learn the basics of setting up your account and navigating the dashboard.", icon: Rocket, color: "bg-blue-50 text-blue-600" },
  { title: "Shipments Guide", description: "How to create, track, and manage shipments from pickup to delivery.", icon: Package, color: "bg-green-50 text-green-600" },
  { title: "Fleet Management", description: "Add vehicles, assign drivers, schedule maintenance, and monitor routes.", icon: Truck, color: "bg-orange-50 text-orange-600" },
  { title: "Billing & Invoices", description: "Manage billing details, view invoices, and configure payment methods.", icon: CreditCard, color: "bg-purple-50 text-purple-600" },
  { title: "API Documentation", description: "Integrate CargoMax with your systems using our REST API endpoints.", icon: Code, color: "bg-gray-100 text-gray-700" },
  { title: "FAQs", description: "Answers to the most commonly asked questions about the platform.", icon: HelpCircle, color: "bg-yellow-50 text-yellow-600" },
];

export default function HelpPage() {
  return (
    <div className="p-6 max-w-5xl mx-auto">
      <h1 className="text-2xl font-bold mb-2">Help Center</h1>
      <p className="text-gray-500 mb-6">Find answers, guides, and resources to help you get the most out of CargoMax.</p>

      <div className="relative mb-8">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input
          type="text"
          placeholder="Search for help articles, guides, FAQs..."
          className="w-full pl-10 pr-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none text-sm"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {cards.map((card) => (
          <div
            key={card.title}
            className="bg-white border rounded-lg p-5 hover:shadow-md transition-shadow cursor-pointer group"
          >
            <div className={`w-10 h-10 rounded-lg ${card.color} flex items-center justify-center mb-3`}>
              <card.icon className="w-5 h-5" />
            </div>
            <h3 className="font-semibold text-sm mb-1 group-hover:text-blue-600 transition-colors">
              {card.title}
            </h3>
            <p className="text-xs text-gray-500 leading-relaxed">{card.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
