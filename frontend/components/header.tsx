"use client";

import { useState } from "react";
import { Search, Sun, Moon, Bell, Globe, LogOut } from "lucide-react";
import { useAuth } from "@/lib/auth";

export default function Header() {
  const [darkMode, setDarkMode] = useState(false);
  const { user, logout } = useAuth();

  const initials = user
    ? `${(user.firstName || "")[0] || ""}${(user.lastName || "")[0] || ""}`.toUpperCase() || user.email[0].toUpperCase()
    : "?";

  return (
    <header className="sticky top-0 w-full h-16 bg-white border-b flex items-center justify-between px-6 z-50">
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
        <input
          type="text"
          placeholder="Search shipments, clients, orders..."
          className="rounded-md border bg-gray-50 w-96 h-10 pl-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      <div className="flex items-center gap-4">
        <button onClick={() => setDarkMode(!darkMode)} className="rounded-md hover:bg-gray-100 p-2 transition">
          {darkMode ? <Moon className="h-5 w-5 text-gray-600" /> : <Sun className="h-5 w-5 text-gray-600" />}
        </button>

        <button className="relative rounded-md hover:bg-gray-100 p-2 transition">
          <Bell className="h-5 w-5 text-gray-600" />
          <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">3</span>
        </button>

        <button className="rounded-md hover:bg-gray-100 p-2 transition">
          <Globe className="h-5 w-5 text-gray-600" />
        </button>

        <div className="flex items-center gap-2">
          <div className="rounded-full bg-blue-500 w-8 h-8 flex items-center justify-center text-white text-sm font-medium">
            {initials}
          </div>
          <div className="text-sm">
            <span className="font-medium text-gray-700">{user?.firstName || user?.email}</span>
            <span className="text-gray-400 ml-1 text-xs">{user?.role}</span>
          </div>
          <button onClick={logout} className="rounded-md hover:bg-gray-100 p-2 transition ml-1" title="Sign out">
            <LogOut className="h-4 w-4 text-gray-500" />
          </button>
        </div>
      </div>
    </header>
  );
}
