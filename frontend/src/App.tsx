import './App.css';

export function App() {
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
    </main>
  );
}
