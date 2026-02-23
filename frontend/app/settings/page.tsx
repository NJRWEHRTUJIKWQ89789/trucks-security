"use client";

import { useEffect, useState } from "react";
import { Save, Moon, Sun, Bell, Mail, MessageSquare, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Setting {
  id: string;
  key: string;
  value: string;
  category: string;
  updatedAt: string;
}

const TIMEZONE_OPTIONS = [
  { value: "America/Los_Angeles", label: "UTC-08:00 Pacific" },
  { value: "America/Chicago", label: "UTC-06:00 Central" },
  { value: "America/New_York", label: "UTC-05:00 Eastern" },
  { value: "Europe/London", label: "UTC+00:00 GMT" },
  { value: "Europe/Berlin", label: "UTC+01:00 CET" },
  { value: "Asia/Kolkata", label: "UTC+05:30 IST" },
  { value: "Asia/Shanghai", label: "UTC+08:00 CST" },
  { value: "Asia/Tokyo", label: "UTC+09:00 JST" },
];

const CURRENCY_OPTIONS = [
  { value: "USD", label: "USD ($)" },
  { value: "EUR", label: "EUR (€)" },
  { value: "GBP", label: "GBP (£)" },
  { value: "INR", label: "INR (₹)" },
  { value: "JPY", label: "JPY (¥)" },
];

const DATE_FORMAT_OPTIONS = [
  { value: "MM/DD/YYYY", label: "MM/DD/YYYY" },
  { value: "DD/MM/YYYY", label: "DD/MM/YYYY" },
  { value: "YYYY-MM-DD", label: "YYYY-MM-DD" },
];

export default function SettingsPage() {
  const [settings, setSettings] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [darkMode, setDarkMode] = useState(false);

  useEffect(() => {
    gql<{ settings: Setting[] }>(`{ settings { id key value category updatedAt } }`)
      .then((d) => {
        const map: Record<string, string> = {};
        for (const s of d.settings) {
          map[s.key] = s.value;
        }
        setSettings(map);
        setDarkMode(map["dark_mode"] === "true");
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const updateLocal = (key: string, value: string) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const keysToSave = ["company_name", "timezone", "currency", "date_format", "email_notifications", "sms_notifications", "push_notifications", "dark_mode"];
      for (const key of keysToSave) {
        if (settings[key] !== undefined) {
          await gql(`mutation($key: String!, $value: String!) { updateSetting(key: $key, value: $value) { id key value } }`, { key, value: settings[key] });
        }
      }
    } catch {
      // error handled silently
    } finally {
      setSaving(false);
    }
  };

  const toggleDarkMode = () => {
    const newVal = !darkMode;
    setDarkMode(newVal);
    updateLocal("dark_mode", String(newVal));
  };

  const notifications = {
    email: settings["email_notifications"] === "true",
    sms: settings["sms_notifications"] === "true",
    push: settings["push_notifications"] === "true",
  };

  const toggleNotification = (key: "email" | "sms" | "push") => {
    const settingKey = `${key}_notifications`;
    const newVal = settings[settingKey] === "true" ? "false" : "true";
    updateLocal(settingKey, newVal);
  };

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">General Settings</h1>

      <div className="bg-white rounded-lg border p-6 space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Company Name</label>
          <input
            type="text"
            value={settings["company_name"] || ""}
            onChange={(e) => updateLocal("company_name", e.target.value)}
            className="w-full border rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Timezone</label>
            <select
              value={settings["timezone"] || "America/Chicago"}
              onChange={(e) => updateLocal("timezone", e.target.value)}
              className="w-full border rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
            >
              {TIMEZONE_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Currency</label>
            <select
              value={settings["currency"] || "USD"}
              onChange={(e) => updateLocal("currency", e.target.value)}
              className="w-full border rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
            >
              {CURRENCY_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Date Format</label>
            <select
              value={settings["date_format"] || "MM/DD/YYYY"}
              onChange={(e) => updateLocal("date_format", e.target.value)}
              className="w-full border rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
            >
              {DATE_FORMAT_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="flex items-center justify-between py-3 border-t">
          <div className="flex items-center gap-3">
            {darkMode ? <Moon className="w-5 h-5 text-indigo-500" /> : <Sun className="w-5 h-5 text-yellow-500" />}
            <div>
              <p className="font-medium text-sm">Dark Mode</p>
              <p className="text-xs text-gray-500">Switch between light and dark themes</p>
            </div>
          </div>
          <button
            onClick={toggleDarkMode}
            className={`relative w-11 h-6 rounded-full transition-colors ${darkMode ? "bg-indigo-500" : "bg-gray-300"}`}
          >
            <span
              className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform ${darkMode ? "translate-x-5" : ""}`}
            />
          </button>
        </div>

        <div className="border-t pt-4 space-y-3">
          <p className="font-medium text-sm">Notifications</p>
          {([
            { key: "email" as const, label: "Email Notifications", icon: Mail },
            { key: "sms" as const, label: "SMS Notifications", icon: MessageSquare },
            { key: "push" as const, label: "Push Notifications", icon: Bell },
          ]).map(({ key, label, icon: Icon }) => (
            <div key={key} className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Icon className="w-4 h-4 text-gray-500" />
                <span className="text-sm">{label}</span>
              </div>
              <button
                onClick={() => toggleNotification(key)}
                className={`relative w-11 h-6 rounded-full transition-colors ${notifications[key] ? "bg-blue-500" : "bg-gray-300"}`}
              >
                <span
                  className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform ${notifications[key] ? "translate-x-5" : ""}`}
                />
              </button>
            </div>
          ))}
        </div>

        <div className="border-t pt-4 flex justify-end">
          <button
            onClick={handleSave}
            disabled={saving}
            className="flex items-center gap-2 bg-blue-600 text-white px-5 py-2 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : <Save className="w-4 h-4" />}
            {saving ? "Saving..." : "Save Changes"}
          </button>
        </div>
      </div>
    </div>
  );
}
