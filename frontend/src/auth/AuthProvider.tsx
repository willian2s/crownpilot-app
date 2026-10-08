import {
  GoogleAuthProvider,
  getIdToken,
  onIdTokenChanged,
  signInWithPopup,
  signInWithRedirect,
  signOut as firebaseSignOut,
  type Auth,
  type User,
} from 'firebase/auth';
import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from 'react';
import { getFirebaseAuth } from './firebase';
import { AuthContext, type AuthContextValue } from './AuthContext';

const configurationError = 'Google Sign-In is not configured for this environment.';

function errorCode(error: unknown): string | undefined {
  if (typeof error !== 'object' || error === null || !('code' in error)) {
    return undefined;
  }

  const code = error.code;
  return typeof code === 'string' ? code : undefined;
}

function publicAuthError(error: unknown): string {
  return errorCode(error) === 'auth/popup-closed-by-user'
    ? 'Google Sign-In was cancelled.'
    : 'Google Sign-In failed. Try again.';
}

function requireAuth(): Auth {
  const auth = getFirebaseAuth();
  if (!auth) {
    throw new Error(configurationError);
  }

  return auth;
}

export function AuthProvider({ children }: PropsWithChildren) {
  const auth = getFirebaseAuth();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(auth !== null);
  const [error, setError] = useState<string | null>(auth ? null : configurationError);

  useEffect(() => {
    if (!auth) {
      return undefined;
    }

    let mounted = true;
    const unsubscribe = onIdTokenChanged(
      auth,
      (currentUser) => {
        if (!mounted) {
          return;
        }

        setUser(currentUser);
        setLoading(false);
      },
      (authError) => {
        if (!mounted) {
          return;
        }

        setUser(null);
        setLoading(false);
        setError(publicAuthError(authError));
      },
    );

    return () => {
      mounted = false;
      unsubscribe();
    };
  }, [auth]);

  const signInWithGoogle = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const auth = requireAuth();
      const provider = new GoogleAuthProvider();
      try {
        await signInWithPopup(auth, provider);
      } catch (popupError) {
        if (errorCode(popupError) !== 'auth/popup-blocked') {
          throw popupError;
        }

        await signInWithRedirect(auth, provider);
      }
    } catch (signInError) {
      setError(errorCode(signInError) === undefined && signInError instanceof Error &&
          signInError.message === configurationError
        ? configurationError
        : publicAuthError(signInError));
    } finally {
      setLoading(false);
    }
  }, []);

  const signOut = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      await firebaseSignOut(requireAuth());
    } catch (signOutError) {
      setError(publicAuthError(signOutError));
    } finally {
      setLoading(false);
    }
  }, []);

  const getCurrentIdToken = useCallback(async (forceRefresh = false) => {
    const auth = getFirebaseAuth();
    const currentUser = auth?.currentUser;
    return currentUser ? getIdToken(currentUser, forceRefresh) : null;
  }, []);

  const value = useMemo<AuthContextValue>(() => ({
    user,
    loading,
    error,
    signInWithGoogle,
    signOut,
    getIdToken: getCurrentIdToken,
  }), [error, getCurrentIdToken, loading, signInWithGoogle, signOut, user]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
