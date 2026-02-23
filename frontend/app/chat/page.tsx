"use client";

import { useState } from "react";
import { Send, Paperclip, Smile, Search, MoreVertical } from "lucide-react";

const contacts = [
  { name: "Sarah Chen", role: "Fleet Manager", online: true, lastMsg: "Sounds good, I'll update the route." },
  { name: "Mike Rodriguez", role: "Dispatcher", online: true, lastMsg: "ETA for SHP-4521 is 3:30 PM." },
  { name: "Laura Kim", role: "Warehouse Ops", online: false, lastMsg: "Denver Hub is fully operational now." },
  { name: "James Wilson", role: "Mechanic Lead", online: true, lastMsg: "VH-208 service is done." },
  { name: "David Park", role: "Driver", online: false, lastMsg: "Arrived at the pickup point." },
  { name: "Emily Torres", role: "Billing", online: false, lastMsg: "Invoice sent to client." },
  { name: "Ryan Foster", role: "Driver", online: true, lastMsg: "On my way to the depot." },
];

const messages = [
  { from: "them", text: "Hey, quick update on the Denver route. We rerouted two trucks to avoid the I-70 closure.", time: "2:14 PM" },
  { from: "me", text: "Got it, thanks for the heads up. What's the new ETA for the deliveries?", time: "2:16 PM" },
  { from: "them", text: "Truck A should arrive by 4:30 PM, Truck B around 5:15 PM. Both within the delivery window.", time: "2:17 PM" },
  { from: "me", text: "Perfect. Can you send the updated route map to the warehouse team?", time: "2:19 PM" },
  { from: "them", text: "Already done! Laura confirmed she got it. They're prepping the unloading bays.", time: "2:20 PM" },
  { from: "me", text: "Great work. Let me know if anything else comes up.", time: "2:22 PM" },
  { from: "them", text: "Will do. Also, VH-312 is running low on fuel — should we reroute to the depot first?", time: "2:25 PM" },
  { from: "me", text: "Yes, have Ryan stop at the depot. We can't risk a breakdown on the highway.", time: "2:26 PM" },
];

export default function ChatPage() {
  const [input, setInput] = useState("");
  const [activeContact, setActiveContact] = useState(contacts[0]);

  return (
    <div className="flex h-[calc(100vh-4rem)]">
      {/* Contacts */}
      <div className="w-72 lg:w-80 border-r bg-white flex flex-col shrink-0">
        <div className="p-3 border-b">
          <div className="relative">
            <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search contacts..."
              className="w-full pl-8 pr-3 py-2 bg-gray-100 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
        <div className="flex-1 overflow-y-auto">
          {contacts.map((c) => (
            <div
              key={c.name}
              onClick={() => setActiveContact(c)}
              className={`flex items-center gap-3 px-3 py-3 cursor-pointer transition-colors ${
                activeContact.name === c.name ? "bg-blue-50" : "hover:bg-gray-50"
              }`}
            >
              <div className="relative">
                <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-sm font-medium text-gray-600">
                  {c.name.split(" ").map((n) => n[0]).join("")}
                </div>
                {c.online && (
                  <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-white rounded-full" />
                )}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium truncate">{c.name}</p>
                <p className="text-xs text-gray-400 truncate">{c.lastMsg}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Chat Area */}
      <div className="flex-1 flex flex-col bg-gray-50">
        {/* Header */}
        <div className="flex items-center gap-3 px-4 py-3 bg-white border-b">
          <div className="w-9 h-9 rounded-full bg-gray-200 flex items-center justify-center text-sm font-medium text-gray-600">
            {activeContact.name.split(" ").map((n) => n[0]).join("")}
          </div>
          <div>
            <p className="text-sm font-semibold">{activeContact.name}</p>
            <p className="text-xs text-gray-400">{activeContact.role} {activeContact.online ? "· Online" : ""}</p>
          </div>
          <button className="ml-auto p-2 hover:bg-gray-100 rounded-lg">
            <MoreVertical className="w-4 h-4 text-gray-500" />
          </button>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {messages.map((msg, i) => (
            <div key={i} className={`flex ${msg.from === "me" ? "justify-end" : "justify-start"}`}>
              <div
                className={`max-w-[70%] px-3 py-2 rounded-2xl text-sm ${
                  msg.from === "me"
                    ? "bg-blue-600 text-white rounded-br-md"
                    : "bg-white text-gray-800 border rounded-bl-md"
                }`}
              >
                <p>{msg.text}</p>
                <p className={`text-[10px] mt-1 ${msg.from === "me" ? "text-blue-200" : "text-gray-400"}`}>
                  {msg.time}
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* Input */}
        <div className="p-3 bg-white border-t">
          <div className="flex items-center gap-2">
            <button className="p-2 hover:bg-gray-100 rounded-lg">
              <Paperclip className="w-4 h-4 text-gray-500" />
            </button>
            <button className="p-2 hover:bg-gray-100 rounded-lg">
              <Smile className="w-4 h-4 text-gray-500" />
            </button>
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Type a message..."
              className="flex-1 px-3 py-2 bg-gray-100 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button className="p-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
              <Send className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
