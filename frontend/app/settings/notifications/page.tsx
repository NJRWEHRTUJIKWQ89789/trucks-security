"use client";

import { useEffect, useState } from "react";
import { Bell, Mail, MessageSquare, Smartphone, Save, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface NotificationPreference {
  id: string;
  eventType: string;
  emailEnabled: boolean;
  smsEnabled: boolean;
  pushEnabled: boolean;
}

// Map event types from the backend to display labels and descriptions
const EVENT_TYPE_META: Record<string, { label: string; description: string }> = {
  shipment_update: { label: "Shipment Updates", description: "When shipments change status, are delayed, or delivered" },
  order_update: { label: "Order Updates", description: "New orders and order status changes" },
  fleet_alert: { label: "Fleet Alerts", description: "Vehicle maintenance, fuel, and driver assignment notifications" },
  system_notification: { label: "System Notifications", description: "Platform updates, reports ready, and system announcements" },
};

type Channel = "emailEnabled" | "smsEnabled" | "pushEnabled";

export default function NotificationsPage() {
  const [preferences, setPreferences] = useState<NotificationPreference[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    gql<{ notificationPreferences: NotificationPreference[] }>(
      `{ notificationPreferences { id eventType emailEnabled smsEnabled pushEnabled } }`
    )
      .then((d) => setPreferences(d.notificationPreferences))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const toggle = (eventType: string, channel: Channel) => {
    setPreferences((prev) =>
      prev.map((p) =>
        p.eventType === eventType ? { ...p, [channel]: !p[channel] } : p
      )
    );
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      for (const pref of preferences) {
        await gql(
          `mutation($input: NotificationPrefInput!) { updateNotificationPreference(input: $input) { id eventType emailEnabled smsEnabled pushEnabled } }`,
          {
            input: {
              eventType: pref.eventType,
              emailEnabled: pref.emailEnabled,
              smsEnabled: pref.smsEnabled,
              pushEnabled: pref.pushEnabled,
            },
          }
        );
      }
    } catch {
      // error handled silently
    } finally {
      setSaving(false);
    }
  };

  const Toggle = ({ on, onToggle }: { on: boolean; onToggle: () => void }) => (
    <button
      onClick={onToggle}
      className={`relative w-10 h-5 rounded-full transition-colors ${on ? "bg-blue-500" : "bg-gray-300"}`}
    >
      <span
        className={`absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full transition-transform ${on ? "translate-x-5" : ""}`}
      />
    </button>
  );

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Bell className="w-6 h-6 text-blue-600" />
          <h1 className="text-2xl font-bold">Notification Preferences</h1>
        </div>
        <button
          onClick={handleSave}
          disabled={saving}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
        >
          {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : <Save className="w-4 h-4" />}
          {saving ? "Saving..." : "Save"}
        </button>
      </div>

      <div className="bg-white rounded-lg border">
        <div className="grid grid-cols-[1fr_80px_80px_80px] gap-4 px-4 py-3 bg-gray-50 border-b text-sm font-semibold text-gray-600">
          <span>Event</span>
          <span className="flex items-center justify-center gap-1"><Mail className="w-4 h-4" /> Email</span>
          <span className="flex items-center justify-center gap-1"><MessageSquare className="w-4 h-4" /> SMS</span>
          <span className="flex items-center justify-center gap-1"><Smartphone className="w-4 h-4" /> Push</span>
        </div>

        {preferences.length === 0 && (
          <div className="px-4 py-8 text-center text-gray-500 text-sm">
            No notification preferences configured.
          </div>
        )}

        {preferences.map((pref) => {
          const meta = EVENT_TYPE_META[pref.eventType] || {
            label: pref.eventType.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase()),
            description: "",
          };
          return (
            <div
              key={pref.eventType}
              className="grid grid-cols-[1fr_80px_80px_80px] gap-4 px-4 py-4 border-b last:border-0 items-center hover:bg-gray-50"
            >
              <div>
                <p className="font-medium text-sm">{meta.label}</p>
                <p className="text-xs text-gray-500">{meta.description}</p>
              </div>
              <div className="flex justify-center">
                <Toggle on={pref.emailEnabled} onToggle={() => toggle(pref.eventType, "emailEnabled")} />
              </div>
              <div className="flex justify-center">
                <Toggle on={pref.smsEnabled} onToggle={() => toggle(pref.eventType, "smsEnabled")} />
              </div>
              <div className="flex justify-center">
                <Toggle on={pref.pushEnabled} onToggle={() => toggle(pref.eventType, "pushEnabled")} />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
