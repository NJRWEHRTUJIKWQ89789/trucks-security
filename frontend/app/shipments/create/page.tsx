"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Loader2, CheckCircle, AlertCircle } from "lucide-react";
import { gql } from "@/lib/graphql";

interface ShipmentInput {
  trackingNumber: string;
  origin?: string;
  destination?: string;
  carrier?: string;
  weight?: number;
  dimensions?: string;
  estimatedDelivery?: string;
  customerName?: string;
  customerEmail?: string;
  notes?: string;
}

interface CreatedShipment {
  id: string;
  trackingNumber: string;
  status: string;
}

export default function CreateShipmentPage() {
  const router = useRouter();
  const [shippingMethod, setShippingMethod] = useState("standard");
  const [insurance, setInsurance] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<CreatedShipment | null>(null);

  // Form field state
  const [senderName, setSenderName] = useState("");
  const [senderPhone, setSenderPhone] = useState("");
  const [senderEmail, setSenderEmail] = useState("");
  const [senderAddress, setSenderAddress] = useState("");
  const [receiverName, setReceiverName] = useState("");
  const [receiverPhone, setReceiverPhone] = useState("");
  const [receiverEmail, setReceiverEmail] = useState("");
  const [receiverAddress, setReceiverAddress] = useState("");
  const [weight, setWeight] = useState("");
  const [dimensions, setDimensions] = useState("");
  const [packageType, setPackageType] = useState("Box");
  const [description, setDescription] = useState("");

  const generateTrackingNumber = () => {
    const rand = Math.random().toString(36).substring(2, 8).toUpperCase();
    return `SH-${new Date().getFullYear()}-${rand}`;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(null);

    const trackingNumber = generateTrackingNumber();

    const notes = [
      description && `Contents: ${description}`,
      packageType && `Package type: ${packageType}`,
      shippingMethod && `Shipping method: ${shippingMethod}`,
      insurance && "Insured",
      senderPhone && `Sender phone: ${senderPhone}`,
      receiverPhone && `Receiver phone: ${receiverPhone}`,
    ]
      .filter(Boolean)
      .join("; ");

    const input: ShipmentInput = {
      trackingNumber,
      origin: senderAddress || undefined,
      destination: receiverAddress || undefined,
      carrier: shippingMethod === "overnight" ? "Express Air" : shippingMethod === "express" ? "Express Freight" : "Standard Ground",
      weight: weight ? parseFloat(weight) : undefined,
      dimensions: dimensions || undefined,
      customerName: receiverName || undefined,
      customerEmail: receiverEmail || undefined,
      notes: notes || undefined,
    };

    try {
      const data = await gql<{ createShipment: CreatedShipment }>(
        `mutation($input:ShipmentInput!){createShipment(input:$input){id trackingNumber status}}`,
        { input }
      );
      setSuccess(data.createShipment);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Failed to create shipment";
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  if (success) {
    return (
      <div className="p-8 max-w-4xl mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8 text-center">
          <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Shipment Created</h1>
          <p className="text-gray-600 mb-6">
            Your shipment has been created with tracking number{" "}
            <span className="font-semibold text-blue-600">{success.trackingNumber}</span>.
          </p>
          <div className="flex items-center justify-center gap-3">
            <button
              onClick={() => router.push(`/shipments/track`)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Track Shipment
            </button>
            <button
              onClick={() => {
                setSuccess(null);
                setSenderName("");
                setSenderPhone("");
                setSenderEmail("");
                setSenderAddress("");
                setReceiverName("");
                setReceiverPhone("");
                setReceiverEmail("");
                setReceiverAddress("");
                setWeight("");
                setDimensions("");
                setPackageType("Box");
                setDescription("");
                setShippingMethod("standard");
                setInsurance(false);
              }}
              className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Create Another
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Create New Shipment</h1>

      {error && (
        <div className="flex items-center gap-3 bg-red-50 border border-red-200 rounded-xl p-4 text-red-700 mb-6">
          <AlertCircle className="w-5 h-5 flex-shrink-0" />
          <p className="text-sm">{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Sender Info */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Sender Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input
                type="text"
                value={senderName}
                onChange={(e) => setSenderName(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="John Doe"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input
                type="tel"
                value={senderPhone}
                onChange={(e) => setSenderPhone(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="+1 (555) 000-0000"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                type="email"
                value={senderEmail}
                onChange={(e) => setSenderEmail(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="john@example.com"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Address</label>
              <input
                type="text"
                value={senderAddress}
                onChange={(e) => setSenderAddress(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="123 Main St, New York, NY"
              />
            </div>
          </div>
        </div>

        {/* Receiver Info */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Receiver Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              <input
                type="text"
                value={receiverName}
                onChange={(e) => setReceiverName(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Jane Smith"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input
                type="tel"
                value={receiverPhone}
                onChange={(e) => setReceiverPhone(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="+1 (555) 000-0000"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                type="email"
                value={receiverEmail}
                onChange={(e) => setReceiverEmail(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="jane@example.com"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Address</label>
              <input
                type="text"
                value={receiverAddress}
                onChange={(e) => setReceiverAddress(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="456 Oak Ave, Los Angeles, CA"
              />
            </div>
          </div>
        </div>

        {/* Package Details */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Package Details</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Weight (kg)</label>
              <input
                type="number"
                step="0.1"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="0.0"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Dimensions (L x W x H cm)</label>
              <input
                type="text"
                value={dimensions}
                onChange={(e) => setDimensions(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="30 x 20 x 15"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Package Type</label>
              <select
                value={packageType}
                onChange={(e) => setPackageType(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option>Box</option>
                <option>Envelope</option>
                <option>Pallet</option>
                <option>Crate</option>
                <option>Tube</option>
              </select>
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea
                rows={3}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Describe the package contents..."
              />
            </div>
          </div>
        </div>

        {/* Shipping Options */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Shipping Options</h2>
          <div className="space-y-3 mb-4">
            {[
              { value: "standard", label: "Standard", desc: "5-7 business days", price: "$12.99" },
              { value: "express", label: "Express", desc: "2-3 business days", price: "$24.99" },
              { value: "overnight", label: "Overnight", desc: "Next business day", price: "$49.99" },
            ].map((option) => (
              <label
                key={option.value}
                className={`flex items-center justify-between p-4 border rounded-lg cursor-pointer transition-colors ${
                  shippingMethod === option.value ? "border-blue-500 bg-blue-50" : "border-gray-200 hover:bg-gray-50"
                }`}
              >
                <div className="flex items-center gap-3">
                  <input
                    type="radio"
                    name="shipping"
                    value={option.value}
                    checked={shippingMethod === option.value}
                    onChange={(e) => setShippingMethod(e.target.value)}
                    className="w-4 h-4 text-blue-600"
                  />
                  <div>
                    <p className="text-sm font-medium text-gray-900">{option.label}</p>
                    <p className="text-xs text-gray-500">{option.desc}</p>
                  </div>
                </div>
                <p className="text-sm font-semibold text-gray-900">{option.price}</p>
              </label>
            ))}
          </div>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={insurance}
              onChange={(e) => setInsurance(e.target.checked)}
              className="w-4 h-4 text-blue-600 rounded"
            />
            <span className="text-sm text-gray-700">Add shipping insurance (+$4.99)</span>
          </label>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-3 justify-end">
          <button
            type="button"
            onClick={() => router.push("/shipments")}
            className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={submitting}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 flex items-center gap-2"
          >
            {submitting && <Loader2 className="w-4 h-4 animate-spin" />}
            {submitting ? "Creating..." : "Create Shipment"}
          </button>
        </div>
      </form>
    </div>
  );
}
