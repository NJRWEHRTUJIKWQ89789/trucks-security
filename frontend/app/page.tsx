import StatsCards from "@/components/stats-cards";
import ShipmentTrends from "@/components/shipment-trends";
import FleetStatus from "@/components/fleet-status";
import RecentActivity from "@/components/recent-activity";
import ActiveDeliveries from "@/components/active-deliveries";
import QuickActions from "@/components/quick-actions";
import Alerts from "@/components/alerts";
import PerformanceHighlights from "@/components/performance-highlights";
import TopRoutes from "@/components/top-routes";

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Dashboard Overview</h1>
          <p className="text-gray-500">
            Welcome back! Here&apos;s what&apos;s happening with your logistics operations.
          </p>
        </div>
        <div className="flex gap-3">
          <button className="bg-blue-600 text-white rounded-md px-4 py-2 text-sm font-medium hover:bg-blue-700">
            Add Vehicle
          </button>
          <button className="bg-blue-600 text-white rounded-md px-4 py-2 text-sm font-medium hover:bg-blue-700">
            New Shipment
          </button>
        </div>
      </div>
      <StatsCards />
      <div className="grid grid-cols-3 gap-6">
        <div className="col-span-2">
          <ShipmentTrends />
        </div>
        <div className="col-span-1">
          <FleetStatus />
        </div>
      </div>
      <div className="grid grid-cols-2 gap-6">
        <RecentActivity />
        <ActiveDeliveries />
      </div>
      <QuickActions />
      <div className="grid grid-cols-3 gap-6">
        <PerformanceHighlights />
        <TopRoutes />
        <Alerts />
      </div>
    </div>
  );
}
