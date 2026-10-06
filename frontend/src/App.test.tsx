import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { App } from './App';

describe('App', () => {
  it('renders bootstrap status without initializing provider integrations', () => {
    render(<App />);

    expect(screen.getByRole('heading', { name: 'CrownPilot' })).toBeVisible();
    expect(screen.getByRole('status')).toHaveTextContent('Bootstrap status: ready');
    expect(screen.getByText(/Identity, persistence, and player-link flows/)).toBeVisible();
  });
});
