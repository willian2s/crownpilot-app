import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { User } from 'firebase/auth';
import { AuthProvider } from './AuthProvider';
import { useAuth } from './useAuth';

const mocks = vi.hoisted(() => ({
  auth: { currentUser: null },
  currentUser: null as User | null,
  idTokenChanged: null as ((user: User | null) => void) | null,
  idTokenError: null as ((error: unknown) => void) | null,
  signOut: vi.fn(),
}));

vi.mock('./firebase', () => ({
  getFirebaseAuth: () => mocks.auth,
}));

vi.mock('firebase/auth', () => ({
  GoogleAuthProvider: class GoogleAuthProvider {},
  getIdToken: vi.fn(async () => 'firebase-id-token'),
  onIdTokenChanged: vi.fn((_auth, next, error) => {
    mocks.idTokenChanged = next;
    mocks.idTokenError = error;
    next(mocks.currentUser);
    return vi.fn();
  }),
  signInWithPopup: vi.fn(),
  signInWithRedirect: vi.fn(),
  signOut: mocks.signOut,
}));

function AuthProbe() {
  const { error, loading, signOut, user } = useAuth();

  return (
    <>
      <output>{loading ? 'loading' : user ? 'signed-in' : 'signed-out'}</output>
      <button type="button" onClick={() => void signOut()}>Sign out</button>
      {error ? <div role="alert">{error}</div> : null}
    </>
  );
}

describe('AuthProvider', () => {
  afterEach(() => {
    cleanup();
  });

  beforeEach(() => {
    mocks.currentUser = { uid: 'firebase-uid' } as User;
    mocks.idTokenChanged = null;
    mocks.idTokenError = null;
    mocks.signOut.mockReset().mockResolvedValue(undefined);
  });

  it('tracks Firebase token lifecycle and logout through SDK', async () => {
    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => expect(screen.getByText('signed-in')).toBeVisible());
    fireEvent.click(screen.getByRole('button', { name: 'Sign out' }));

    expect(mocks.signOut).toHaveBeenCalledWith(mocks.auth);
    mocks.idTokenChanged?.(null);
    await waitFor(() => expect(screen.getByText('signed-out')).toBeVisible());
  });

  it('clears session and exposes generic error when token refresh fails', async () => {
    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    );

    await waitFor(() => expect(screen.getByText('signed-in')).toBeVisible());
    mocks.idTokenError?.(new Error('provider failure'));

    await waitFor(() => {
      expect(screen.getByText('signed-out')).toBeVisible();
      expect(screen.getByRole('alert')).toHaveTextContent('Google Sign-In failed. Try again.');
    });
  });
});
