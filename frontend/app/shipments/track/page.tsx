"use client";

import { useState } from "react";
import { Search, Check, MapPin, Truck, Package, CircleDot, Loader2, AlertCircle } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Shipment {
  id: string;
  trackingNumber: string;
  origin: string;
  destination: string;
  status: string;
  carrier: string;
  weight: number;
  dimensions: string;
  estimatedDelivery: string;
  actualDelivery: string;
  customerName: string;
  customerEmail: string;
  notes: string;
  createdAt: string;
}

const allSteps = [
  { label: "Order Placed", key: "pending" },
  { label: "Picked Up", key: "picked_up" },
  { label: "In Transit", key: "in_transit" },
  { label: "Out for Delivery", key: "out_for_delivery" },
  { label: "Delivered", key: "delivered" },
];

function getStepStatus(stepIndex: number, currentStepIndex: number): "completed" | "current" | "upcoming" {
  if (stepIndex < currentStepIndex) return "completed";
  if (stepIndex === currentStepIndex) return "current";
  return "upcoming";
}

function getCurrentStepIndex(status: string): number {
  const statusMap: Record<string, number> = {
    pending: 0,
    picked_up: 1,
    in_transit: 2,
    out_for_delivery: 3,
    delivered: 4,
  };
  return statusMap[status] ?? 0;
}

function formatStatus(status: string): string {
  return status
    .split("_")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

function formatDate(dateStr: string | null | undefined): string {
  if (!dateStr) return "-";
  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return dateStr;
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
}

const statusBadgeColors: Record<string, string> = {
  pending: "bg-amber-100 text-amber-700",
  picked_up: "bg-purple-100 text-purple-700",
  in_transit: "bg-blue-100 text-blue-700",
  out_for_delivery: "bg-indigo-100 text-indigo-700",
  delivered: "bg-green-100 text-green-700",
  delayed: "bg-red-100 text-red-700",
};

export default function TrackShipmentPage() {
  const [trackingNumber, setTrackingNumber] = useState("");
  const [shipment, setShipment] = useState<Shipment | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    const tn = trackingNumber.trim();
    if (!tn) return;

    setLoading(true);
    setError(null);
    setShipment(null);

    try {
      const data = await gql<{ trackShipment: Shipment }>(
        `query($trackingNumber:String!){trackShipment(trackingNumber:$trackingNumber){id trackingNumber origin destination status carrier weight dimensions estimatedDelivery actualDelivery customerName customerEmail notes createdAt}}`,
        { trackingNumber: tn }
      );
      setShipment(data.trackShipment);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Failed to find shipment";
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  const currentStepIndex = shipment ? getCurrentStepIndex(shipment.status) : 0;
  const isDelayed = shipment?.status === "delayed";
  const progressPercent = shipment
    ? isDelayed
      ? 50
      : Math.min(100, (currentStepIndex / (allSteps.length - 1)) * 100)
    : 0;

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Track Shipment</h1>

      <div className="max-w-2xl mx-auto mb-10">
        <div className="relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Enter Tracking Number"
            value={trackingNumber}
            onChange={(e) => setTrackingNumber(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSearch()}
            className="w-full pl-12 pr-32 py-4 text-lg border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={handleSearch}
            disabled={loading}
            className="absolute right-2 top-1/2 -translate-y-1/2 bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {loading ? <Loader2 className="w-5 h-5 animate-spin" /> : "Track"}
          </button>
        </div>
      </div>

      {error && (
        <div className="max-w-4xl mx-auto mb-6">
          <div className="flex items-center gap-3 bg-red-50 border border-red-200 rounded-xl p-4 text-red-700">
            <AlertCircle className="w-5 h-5 flex-shrink-0" />
            <p className="text-sm">{error}</p>
          </div>
        </div>
      )}

      {loading && (
        <div className="flex justify-center py-12">
          <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
        </div>
      )}

      {shipment && !loading && (
        <div className="max-w-4xl mx-auto space-y-6">
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-6">
              <div>
                <p className="text-sm text-gray-500">Tracking Number</p>
                <p className="text-xl font-bold text-gray-900">{shipment.trackingNumber}</p>
              </div>
              <span
                className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium ${
                  statusBadgeColors[shipment.status] || "bg-gray-100 text-gray-700"
                }`}
              >
                <Truck className="w-4 h-4" />
                {formatStatus(shipment.status)}
              </span>
            </div>

            <div className="relative flex items-center justify-between">
              {allSteps.map((step, i) => {
                const status = isDelayed ? (i < 3 ? "completed" : "upcoming") : getStepStatus(i, currentStepIndex);
                return (
                  <div key={step.label} className="flex flex-col items-center z-10 flex-1">
                    <div
                      className={`w-10 h-10 rounded-full flex items-center justify-center ${
                        status === "completed"
                          ? "bg-green-500 text-white"
                          : status === "current"
                          ? "bg-blue-500 text-white animate-pulse"
                          : "bg-gray-200 text-gray-400"
                      }`}
                    >
                      {status === "completed" ? (
                        <Check className="w-5 h-5" />
                      ) : status === "current" ? (
                        <CircleDot className="w-5 h-5" />
                      ) : (
                        <span className="w-3 h-3 rounded-full bg-gray-300" />
                      )}
                    </div>
                    <p
                      className={`mt-2 text-xs font-medium text-center ${
                        status === "completed"
                          ? "text-green-600"
                          : status === "current"
                          ? "text-blue-600"
                          : "text-gray-400"
                      }`}
                    >
                      {step.label}
                    </p>
                  </div>
                );
              })}
              <div className="absolute top-5 left-[10%] right-[10%] h-0.5 bg-gray-200 -z-0">
                <div className="h-full bg-green-500" style={{ width: `${progressPercent}%` }} />
              </div>
            </div>

            {isDelayed && (
              <div className="mt-4 flex items-center gap-2 bg-red-50 border border-red-200 rounded-lg p-3 text-red-700 text-sm">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                This shipment is currently delayed.
              </div>
            )}
          </div>

          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Shipment Details</h2>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
              <div>
                <p className="text-sm text-gray-500">From</p>
                <div className="flex items-center gap-1.5 mt-1">
                  <MapPin className="w-4 h-4 text-gray-400" />
                  <p className="text-sm font-medium text-gray-900">{shipment.origin || "-"}</p>
                </div>
              </div>
              <div>
                <p className="text-sm text-gray-500">To</p>
                <div className="flex items-center gap-1.5 mt-1">
                  <MapPin className="w-4 h-4 text-gray-400" />
                  <p className="text-sm font-medium text-gray-900">{shipment.destination || "-"}</p>
                </div>
              </div>
              <div>
                <p className="text-sm text-gray-500">Weight</p>
                <div className="flex items-center gap-1.5 mt-1">
                  <Package className="w-4 h-4 text-gray-400" />
                  <p className="text-sm font-medium text-gray-900">{shipment.weight ? `${shipment.weight} kg` : "-"}</p>
                </div>
              </div>
              <div>
                <p className="text-sm text-gray-500">Carrier</p>
                <p className="text-sm font-medium text-gray-900 mt-1">{shipment.carrier || "-"}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Estimated Delivery</p>
                <p className="text-sm font-medium text-gray-900 mt-1">{formatDate(shipment.estimatedDelivery)}</p>
              </div>
              {shipment.actualDelivery && (
                <div>
                  <p className="text-sm text-gray-500">Actual Delivery</p>
                  <p className="text-sm font-medium text-green-700 mt-1">{formatDate(shipment.actualDelivery)}</p>
                </div>
              )}
              {shipment.dimensions && (
                <div>
                  <p className="text-sm text-gray-500">Dimensions</p>
                  <p className="text-sm font-medium text-gray-900 mt-1">{shipment.dimensions}</p>
                </div>
              )}
              {shipment.customerName && (
                <div>
                  <p className="text-sm text-gray-500">Customer</p>
                  <p className="text-sm font-medium text-gray-900 mt-1">{shipment.customerName}</p>
                </div>
              )}
              {shipment.customerEmail && (
                <div>
                  <p className="text-sm text-gray-500">Customer Email</p>
                  <p className="text-sm font-medium text-gray-900 mt-1">{shipment.customerEmail}</p>
                </div>
              )}
            </div>
            {shipment.notes && (
              <div className="mt-4 pt-4 border-t border-gray-100">
                <p className="text-sm text-gray-500 mb-1">Notes</p>
                <p className="text-sm text-gray-700">{shipment.notes}</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
