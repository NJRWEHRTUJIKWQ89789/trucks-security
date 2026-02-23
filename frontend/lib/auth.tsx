"use client";

import { createContext, useContext, useState, useEffect, ReactNode, useCallback } from "react";
import { gql } from "./graphql";

interface User {
  id: string;
  email: string;
  firstName: string | null;
  lastName: string | null;
  role: string;
  avatarUrl: string | null;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refetch: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  login: async () => {},
  logout: async () => {},
  refetch: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchMe = useCallback(async () => {
    try {
      const data = await gql<{ me: User }>(`{ me { id email firstName lastName role avatarUrl } }`);
      setUser(data.me);
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMe();
  }, [fetchMe]);

  const login = async (email: string, password: string) => {
    const data = await gql<{ login: { user: User; token: string } }>(
      `mutation Login($input: LoginInput!) { login(input: $input) { user { id email firstName lastName role avatarUrl } token } }`,
      { input: { email, password } }
    );
    setUser(data.login.user);
  };

  const logout = async () => {
    try {
      await gql(`mutation { logout }`);
    } catch {
      // ignore
    }
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, refetch: fetchMe }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
