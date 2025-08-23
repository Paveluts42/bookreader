import { create } from "zustand";

interface AuthState {
  token: string | null;
  refreshToken: string | null;
  setTokens: (token: string | null, refreshToken: string | null) => void;
  logout: () => void;
  user?: { userId: string; username: string; isAdmin: boolean };
  setUser: (user: { userId: string; username: string; isAdmin: boolean } | undefined) => void;
}
export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem("token"),
  refreshToken: localStorage.getItem("refreshToken"),

  setTokens: (token, refreshToken) => {
    if (token) {
      localStorage.setItem("token", token);
    } else {
      localStorage.removeItem("token");
    }
    if (refreshToken) {
      localStorage.setItem("refreshToken", refreshToken);
    } else {
      localStorage.removeItem("refreshToken");
    }
    set({ token, refreshToken });
  },
  setUser: (user) => {
    set({ user });
  },

  logout: () => {
    localStorage.removeItem("token");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("user");
    set({ token: null, refreshToken: null });
  },
}));