import './App.css';
import { useAuth } from './auth/useAuth';

export function App() {
  const { error, loading, signInWithGoogle, signOut, user } = useAuth();

  const handleSignIn = async () => {
    await signInWithGoogle();
  };

  const handleSignOut = async () => {
    await signOut();
  };

  return (
    <main className="shell" aria-labelledby="app-title">
      <p className="eyebrow">Application foundation</p>
      <h1 id="app-title">CrownPilot</h1>
      <p className="intro">
        React and Vite are ready. Identity, persistence, and player-link flows are added by later tasks.
      </p>
      <p className="status" role="status">
        Bootstrap status: ready
      </p>
      <section className="identity-card" aria-labelledby="identity-title">
        <h2 id="identity-title">CrownPilot identity</h2>
        {loading ? <p aria-live="polite">Checking Google Sign-In session…</p> : null}
        {user ? (
          <>
            <p aria-live="polite">Signed in with Google.</p>
            <button type="button" onClick={() => void handleSignOut()} disabled={loading}>
              Sign out
            </button>
          </>
        ) : (
          <button type="button" onClick={() => void handleSignIn()} disabled={loading}>
            Continue with Google
          </button>
        )}
        {error ? <p className="error" role="alert">{error}</p> : null}
      </section>
    </main>
  );
}
