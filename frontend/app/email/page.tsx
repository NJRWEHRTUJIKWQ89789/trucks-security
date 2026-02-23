"use client";

import { useState } from "react";
import { Inbox, Star, Paperclip, Reply, Forward, Trash2 } from "lucide-react";

const emails = [
  { id: 1, sender: "Sarah Chen", subject: "Q1 Fleet Report Ready", preview: "Hi, the Q1 fleet utilization report has been compiled and is ready for your review...", time: "10:32 AM", unread: true },
  { id: 2, sender: "Mike Rodriguez", subject: "Re: Route Optimization Update", preview: "Thanks for the feedback. I've adjusted the northern corridor routes as discussed...", time: "9:15 AM", unread: true },
  { id: 3, sender: "Billing Team", subject: "Invoice #INV-2026-0284", preview: "Your invoice for February 2026 has been generated. Total amount: $12,450.00...", time: "8:47 AM", unread: false },
  { id: 4, sender: "James Wilson", subject: "Maintenance Alert: VH-208", preview: "Vehicle VH-208 is due for scheduled maintenance. Please confirm the service date...", time: "Yesterday", unread: true },
  { id: 5, sender: "Laura Kim", subject: "New Warehouse Onboarding", preview: "The Denver Hub setup is complete. Here are the access credentials and operational...", time: "Yesterday", unread: false },
  { id: 6, sender: "System Alerts", subject: "Shipment SHP-4510 Delayed", preview: "Shipment SHP-4510 has been delayed due to weather conditions in the midwest region...", time: "Feb 20", unread: false },
  { id: 7, sender: "David Park", subject: "Driver Training Schedule", preview: "Attached is the updated training schedule for March. Please review the new safety...", time: "Feb 19", unread: false },
  { id: 8, sender: "HR Department", subject: "Policy Update: Remote Work", preview: "Please review the updated remote work policy effective March 1st, 2026...", time: "Feb 18", unread: false },
];

export default function EmailPage() {
  const [selected, setSelected] = useState(emails[0]);

  return (
    <div className="flex h-[calc(100vh-4rem)]">
      {/* Email List */}
      <div className="w-80 lg:w-96 border-r bg-white overflow-y-auto shrink-0">
        <div className="flex items-center gap-2 px-4 py-3 border-b">
          <Inbox className="w-5 h-5 text-blue-600" />
          <h2 className="font-semibold">Inbox</h2>
          <span className="ml-auto text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full font-medium">
            {emails.filter((e) => e.unread).length} new
          </span>
        </div>
        {emails.map((email) => (
          <div
            key={email.id}
            onClick={() => setSelected(email)}
            className={`px-4 py-3 border-b cursor-pointer transition-colors ${
              selected.id === email.id ? "bg-blue-50" : "hover:bg-gray-50"
            }`}
          >
            <div className="flex items-center gap-2 mb-1">
              {email.unread && <span className="w-2 h-2 bg-blue-500 rounded-full shrink-0" />}
              <span className={`text-sm truncate ${email.unread ? "font-semibold" : "text-gray-700"}`}>
                {email.sender}
              </span>
              <span className="ml-auto text-xs text-gray-400 shrink-0">{email.time}</span>
            </div>
            <p className={`text-sm truncate ${email.unread ? "font-medium" : "text-gray-600"}`}>{email.subject}</p>
            <p className="text-xs text-gray-400 truncate mt-0.5">{email.preview}</p>
          </div>
        ))}
      </div>

      {/* Email Content */}
      <div className="flex-1 bg-white overflow-y-auto">
        <div className="">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-xl font-bold mb-1">{selected.subject}</h1>
              <p className="text-sm text-gray-500">
                From <span className="font-medium text-gray-700">{selected.sender}</span> &middot; {selected.time}
              </p>
            </div>
            <div className="flex items-center gap-1">
              <button className="p-2 hover:bg-gray-100 rounded-lg"><Star className="w-4 h-4 text-gray-400" /></button>
              <button className="p-2 hover:bg-gray-100 rounded-lg"><Trash2 className="w-4 h-4 text-gray-400" /></button>
            </div>
          </div>

          <div className="border-t pt-4 text-sm text-gray-700 leading-relaxed space-y-3">
            <p>Hi,</p>
            <p>{selected.preview} We have analyzed the data across all active routes and identified several key patterns that could help improve efficiency.</p>
            <p>Key highlights:</p>
            <ul className="list-disc pl-5 space-y-1">
              <li>Overall fleet utilization increased by 12% compared to Q4 2025</li>
              <li>Fuel costs reduced by 8% through optimized routing</li>
              <li>On-time delivery rate improved to 94.7%</li>
              <li>Three vehicles flagged for upcoming maintenance</li>
            </ul>
            <p>Please review the attached report and let me know if you have any questions or need additional analysis.</p>
            <p>Best regards,<br />{selected.sender}</p>
          </div>

          <div className="flex items-center gap-2 mt-6 pt-4 border-t">
            <button className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50 text-sm">
              <Reply className="w-4 h-4" /> Reply
            </button>
            <button className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50 text-sm">
              <Forward className="w-4 h-4" /> Forward
            </button>
            <button className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50 text-sm">
              <Paperclip className="w-4 h-4" /> Attachments (1)
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
