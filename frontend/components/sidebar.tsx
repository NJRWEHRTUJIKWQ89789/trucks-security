"use client";

import {
  LayoutDashboard, Map, Truck, Package, Search, Plus, AlertTriangle,
  Wrench, Users, Building, BarChart3, RefreshCw, Store, UserPlus,
  MessageSquare, ShoppingCart, Calendar, RotateCcw, XCircle, TrendingUp,
  DollarSign, Gauge, Settings, Shield, Bell, HelpCircle, Phone, Mail,
  MessageCircle, Ticket, FileText, LayoutGrid,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import type { LucideIcon } from "lucide-react";

interface NavItem { label: string; icon: LucideIcon; href: string; }
interface NavSection { heading: string; items: NavItem[]; }

const navSections: NavSection[] = [
  { heading: "Dashboard", items: [
    { label: "Overview", icon: LayoutDashboard, href: "/" },
    { label: "Live Shipment Map", icon: Map, href: "/dashboard/map" },
    { label: "Fleet Status", icon: Truck, href: "/dashboard/fleet-status" },
  ]},
  { heading: "Shipments", items: [
    { label: "All Shipments", icon: Package, href: "/shipments" },
    { label: "Track Shipment", icon: Search, href: "/shipments/track" },
    { label: "Create Shipment", icon: Plus, href: "/shipments/create" },
    { label: "Delayed Shipments", icon: AlertTriangle, href: "/shipments/delayed" },
  ]},
  { heading: "Fleet Management", items: [
    { label: "Vehicle List", icon: Truck, href: "/fleet/vehicles" },
    { label: "Maintenance Logs", icon: Wrench, href: "/fleet/maintenance" },
    { label: "Driver Assignments", icon: Users, href: "/fleet/drivers" },
  ]},
  { heading: "Warehouses", items: [
    { label: "Warehouse Locations", icon: Building, href: "/warehouses" },
    { label: "Inventory Levels", icon: BarChart3, href: "/warehouses/inventory" },
    { label: "Restock Requests", icon: RefreshCw, href: "/warehouses/restock" },
  ]},
  { heading: "Vendors & Clients", items: [
    { label: "Vendor Directory", icon: Store, href: "/vendors" },
    { label: "Add Vendor", icon: UserPlus, href: "/vendors/add" },
    { label: "Clients List", icon: Users, href: "/clients" },
    { label: "Client Feedback", icon: MessageSquare, href: "/clients/feedback" },
  ]},
  { heading: "Orders", items: [
    { label: "All Orders", icon: ShoppingCart, href: "/orders" },
    { label: "Scheduled Deliveries", icon: Calendar, href: "/orders/scheduled" },
    { label: "Returns", icon: RotateCcw, href: "/orders/returns" },
    { label: "Cancellations", icon: XCircle, href: "/orders/cancellations" },
  ]},
  { heading: "Reports", items: [
    { label: "Delivery Performance", icon: TrendingUp, href: "/reports/delivery" },
    { label: "Revenue Analysis", icon: DollarSign, href: "/reports/revenue" },
    { label: "Fleet Efficiency", icon: Gauge, href: "/reports/fleet" },
  ]},
  { heading: "System Tools", items: [
    { label: "Settings", icon: Settings, href: "/settings" },
    { label: "Roles & Permissions", icon: Shield, href: "/settings/roles" },
    { label: "Notifications Setup", icon: Bell, href: "/settings/notifications" },
  ]},
  { heading: "Help & Logs", items: [
    { label: "Help Center", icon: HelpCircle, href: "/help" },
    { label: "Contact", icon: Phone, href: "/contact" },
    { label: "Email", icon: Mail, href: "/email" },
    { label: "Chat", icon: MessageCircle, href: "/chat" },
    { label: "Support Tickets", icon: Ticket, href: "/help/tickets" },
    { label: "Audit Logs", icon: FileText, href: "/help/logs" },
    { label: "Widgets", icon: LayoutGrid, href: "/widgets" },
  ]},
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-white border-r border-gray-200 overflow-y-auto z-40">
      <div className="flex items-center gap-2 px-4 py-5 border-b border-gray-200">
        <Package className="h-6 w-6 text-blue-600" />
        <span className="text-xl font-bold text-gray-900">CargoMax</span>
      </div>
      <nav className="px-2 pb-6">
        {navSections.map((section) => (
          <div key={section.heading}>
            <h3 className="uppercase text-xs font-semibold text-gray-400 tracking-wider px-3 pt-4 pb-2">
              {section.heading}
            </h3>
            <ul className="space-y-0.5">
              {section.items.map((item) => {
                const Icon = item.icon;
                const isActive = pathname === item.href;
                return (
                  <li key={item.label}>
                    <Link href={item.href}
                      className={`flex items-center gap-3 w-full px-3 py-2 rounded-md text-sm transition-colors ${
                        isActive ? "bg-blue-50 text-blue-600 font-medium" : "text-gray-600 hover:bg-gray-100"
                      }`}>
                      <Icon className="h-4 w-4 flex-shrink-0" />
                      <span>{item.label}</span>
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>
    </aside>
  );
}
