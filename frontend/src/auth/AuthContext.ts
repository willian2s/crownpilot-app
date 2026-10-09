import { createContext } from 'react';
import type { User } from 'firebase/auth';

export interface AuthContextValue {
  user: User | null;
  loading: boolean;
  error: string | null;
  signInWithGoogle: () => Promise<void>;
  signOut: () => Promise<void>;
  getIdToken: (forceRefresh?: boolean) => Promise<string | null>;
}

export const defaultAuthContext: AuthContextValue = {
  user: null,
  loading: false,
  error: null,
  signInWithGoogle: async () => undefined,
  signOut: async () => undefined,
  getIdToken: async () => null,
};

export const AuthContext = createContext<AuthContextValue>(defaultAuthContext);
