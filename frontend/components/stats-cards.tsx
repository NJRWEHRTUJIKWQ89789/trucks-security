"use client";

import { useEffect, useState } from "react";
import { gql } from "@/lib/graphql";
import {
  Package, CheckCircle, Clock, DollarSign, Truck, Users, Building, UserCheck,
  ArrowUp, ArrowDown, Loader2,
} from "lucide-react";
import { ReactNode } from "react";

interface DashboardStats {
  activeShipments: number;
  deliveredToday: number;
  pendingOrders: number;
  totalRevenue: number;
  totalVehicles: number;
  activeVehicles: number;
  totalClients: number;
  totalWarehouses: number;
  totalInventory: number;
  totalDrivers: number;
  availableDrivers: number;
}

interface StatCard {
  title: string;
  value: string;
  subtitle: string;
  icon: ReactNode;
  iconBg: string;
  badgeValue: string;
  badgeType: "green" | "red";
  badgeIcon?: "up" | "down";
}

function StatCardItem({ card }: { card: StatCard }) {
  return (
    <div className="bg-white rounded-lg border shadow-sm p-6">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-gray-500">{card.title}</h3>
        <div className={`w-10 h-10 rounded-full ${card.iconBg} flex items-center justify-center`}>{card.icon}</div>
      </div>
      <div className="text-3xl font-bold text-gray-900 mt-2">{card.value}</div>
      <div className="flex items-center justify-between mt-1">
        <span className="text-sm text-gray-500">{card.subtitle}</span>
        <span className={`inline-flex items-center text-xs font-medium px-2 py-0.5 rounded-full ${card.badgeType === "green" ? "bg-green-50 text-green-600" : "bg-red-50 text-red-600"}`}>
          {card.badgeIcon === "up" && <ArrowUp className="w-3 h-3 mr-0.5" />}
          {card.badgeIcon === "down" && <ArrowDown className="w-3 h-3 mr-0.5" />}
          {card.badgeValue}
        </span>
      </div>
    </div>
  );
}

export default function StatsCards() {
  const [stats, setStats] = useState<DashboardStats | null>(null);

  useEffect(() => {
    gql<{ dashboardStats: DashboardStats }>(`{
      dashboardStats {
        activeShipments deliveredToday pendingOrders totalRevenue
        totalVehicles activeVehicles totalClients totalWarehouses totalInventory totalDrivers availableDrivers
      }
    }`).then((d) => setStats(d.dashboardStats)).catch(() => {});
  }, []);

  if (!stats) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  const fleetUtil = stats.totalVehicles > 0 ? Math.round((stats.activeVehicles / stats.totalVehicles) * 100) : 0;
  const driverUtil = stats.totalDrivers > 0 ? Math.round((stats.availableDrivers / stats.totalDrivers) * 100) : 0;

  const row1: StatCard[] = [
    { title: "Active Shipments", value: String(stats.activeShipments), subtitle: "Currently in transit", icon: <Package className="w-5 h-5 text-blue-500" />, iconBg: "bg-blue-50", badgeValue: "Live", badgeType: "green", badgeIcon: "up" },
    { title: "Delivered Today", value: String(stats.deliveredToday), subtitle: "Successful deliveries", icon: <CheckCircle className="w-5 h-5 text-green-500" />, iconBg: "bg-green-50", badgeValue: "Today", badgeType: "green" },
    { title: "Pending Orders", value: String(stats.pendingOrders), subtitle: "Awaiting processing", icon: <Clock className="w-5 h-5 text-amber-500" />, iconBg: "bg-amber-50", badgeValue: "Pending", badgeType: "red", badgeIcon: "down" },
    { title: "Revenue (Total)", value: `$${stats.totalRevenue.toLocaleString("en-US", { minimumFractionDigits: 0 })}`, subtitle: "All time", icon: <DollarSign className="w-5 h-5 text-purple-500" />, iconBg: "bg-purple-50", badgeValue: "Total", badgeType: "green" },
  ];

  const row2: StatCard[] = [
    { title: "Fleet Utilization", value: `${fleetUtil}%`, subtitle: `${stats.activeVehicles}/${stats.totalVehicles} vehicles`, icon: <Truck className="w-5 h-5 text-blue-500" />, iconBg: "bg-blue-50", badgeValue: `${stats.activeVehicles} active`, badgeType: "green" },
    { title: "Active Clients", value: String(stats.totalClients), subtitle: "Total active clients", icon: <Users className="w-5 h-5 text-green-500" />, iconBg: "bg-green-50", badgeValue: "Active", badgeType: "green" },
    { title: "Warehouses", value: String(stats.totalWarehouses), subtitle: `${stats.totalInventory.toLocaleString()} items in stock`, icon: <Building className="w-5 h-5 text-amber-500" />, iconBg: "bg-amber-50", badgeValue: `${stats.totalInventory} items`, badgeType: "green" },
    { title: "Total Drivers", value: String(stats.totalDrivers), subtitle: `${driverUtil}% availability`, icon: <UserCheck className="w-5 h-5 text-indigo-500" />, iconBg: "bg-indigo-50", badgeValue: `${stats.availableDrivers} available`, badgeType: "green" },
  ];

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-4 gap-4">{row1.map((c) => <StatCardItem key={c.title} card={c} />)}</div>
      <div className="grid grid-cols-4 gap-4">{row2.map((c) => <StatCardItem key={c.title} card={c} />)}</div>
    </div>
  );
}
