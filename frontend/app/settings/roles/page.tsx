"use client";

import { useEffect, useState } from "react";
import { Shield, Save, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Role {
  id: string;
  name: string;
  permissions: string;
}

// The permission domains that match the backend JSON structure
const PERMISSION_DOMAINS = ["shipments", "vehicles", "drivers", "warehouses", "orders", "vendors", "clients", "settings", "reports", "dashboard"] as const;

// Display labels for each domain
const DOMAIN_LABELS: Record<string, string> = {
  shipments: "Shipments",
  vehicles: "Vehicles",
  drivers: "Drivers",
  warehouses: "Warehouses",
  orders: "Orders",
  vendors: "Vendors",
  clients: "Clients",
  settings: "Settings",
  reports: "Reports",
  dashboard: "Dashboard",
};

const ROLE_COLORS: Record<string, string> = {
  Admin: "bg-red-500",
  Manager: "bg-orange-500",
  Dispatcher: "bg-blue-500",
  Driver: "bg-green-500",
  Viewer: "bg-gray-400",
};

type ParsedPermissions = Record<string, string[]>;

function parsePermissions(permStr: string): ParsedPermissions {
  try {
    return JSON.parse(permStr) as ParsedPermissions;
  } catch {
    return {};
  }
}

function hasDomainAccess(perms: ParsedPermissions, domain: string): boolean {
  return Array.isArray(perms[domain]) && perms[domain].length > 0;
}

function toggleDomainAccess(perms: ParsedPermissions, domain: string): ParsedPermissions {
  const updated = { ...perms };
  if (hasDomainAccess(updated, domain)) {
    delete updated[domain];
  } else {
    updated[domain] = ["read"];
  }
  return updated;
}

export default function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permState, setPermState] = useState<Record<string, ParsedPermissions>>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    gql<{ roles: Role[] }>(`{ roles { id name permissions } }`)
      .then((d) => {
        setRoles(d.roles);
        const parsed: Record<string, ParsedPermissions> = {};
        for (const role of d.roles) {
          parsed[role.id] = parsePermissions(role.permissions);
        }
        setPermState(parsed);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const toggle = (roleId: string, domain: string) => {
    setPermState((prev) => ({
      ...prev,
      [roleId]: toggleDomainAccess(prev[roleId] || {}, domain),
    }));
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      for (const role of roles) {
        if (role.name === "Admin") continue;
        const perms = permState[role.id];
        if (!perms) continue;
        await gql(
          `mutation($id: String!, $input: RoleInput!) { updateRole(id: $id, input: $input) { id name permissions } }`,
          { id: role.id, input: { name: role.name, permissions: JSON.stringify(perms) } }
        );
      }
    } catch {
      // error handled silently
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Shield className="w-6 h-6 text-blue-600" />
          <h1 className="text-2xl font-bold">Roles & Permissions</h1>
        </div>
        <button
          onClick={handleSave}
          disabled={saving}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
        >
          {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : <Save className="w-4 h-4" />}
          {saving ? "Saving..." : "Save Changes"}
        </button>
      </div>

      <div className="bg-white rounded-lg border overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50 border-b">
              <th className="text-left px-4 py-3 text-sm font-semibold text-gray-600">Role</th>
              {PERMISSION_DOMAINS.map((domain) => (
                <th key={domain} className="px-4 py-3 text-sm font-semibold text-gray-600 text-center">
                  {DOMAIN_LABELS[domain]}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {roles.map((role) => (
              <tr key={role.id} className="border-b last:border-0 hover:bg-gray-50">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <span
                      className={`w-2 h-2 rounded-full ${ROLE_COLORS[role.name] || "bg-gray-400"}`}
                    />
                    <span className="font-medium text-sm">{role.name}</span>
                  </div>
                </td>
                {PERMISSION_DOMAINS.map((domain) => (
                  <td key={domain} className="px-4 py-3 text-center">
                    <input
                      type="checkbox"
                      checked={hasDomainAccess(permState[role.id] || {}, domain)}
                      onChange={() => toggle(role.id, domain)}
                      disabled={role.name === "Admin"}
                      className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
                    />
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <p className="text-xs text-gray-500 mt-3">Admin permissions cannot be modified. All roles inherit read access by default.</p>
    </div>
  );
}
