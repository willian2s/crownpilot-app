package identity

import (
	"errors"
	"testing"
	"time"
)

func TestNewAuthenticatedSubject(t *testing.T) {
	t.Parallel()
	authTime := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		uid     string
		wantErr bool
	}{
		{"verified uid", "firebase-uid-123", false},
		{"empty uid", "", true},
		{"blank uid", "   ", true},
		{"uid with surrounding spaces", " firebase-uid-123 ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAuthenticatedSubject(tt.uid, authTime)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidCredential) {
					t.Fatalf("error = %v, want ErrInvalidCredential", err)
				}
				if !got.IsZero() {
					t.Errorf("subject = %+v, want zero value on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("error = %v", err)
			}
			if got.FirebaseUID() != tt.uid || got.IsZero() {
				t.Errorf("subject = %+v, want uid %q", got, tt.uid)
			}
		})
	}
}

func TestZeroSubjectIsNotAuthenticated(t *testing.T) {
	t.Parallel()
	var s AuthenticatedSubject
	if !s.IsZero() {
		t.Error("zero AuthenticatedSubject must report IsZero")
	}
}

func TestRequireRecentAuthentication(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		authTime time.Time
		wantErr  bool
	}{
		{"just authenticated", now, false},
		{"four minutes ago", now.Add(-4 * time.Minute), false},
		{"exactly at the limit", now.Add(-RecentAuthenticationMaxAge), false},
		{"one second too old", now.Add(-RecentAuthenticationMaxAge - time.Second), true},
		{"one hour ago", now.Add(-time.Hour), true},
		{"slightly in the future (clock skew)", now.Add(30 * time.Second), false},
		{"exactly at the skew tolerance", now.Add(ClockSkewTolerance), false},
		{"beyond the skew tolerance", now.Add(ClockSkewTolerance + time.Second), true},
		{"missing auth_time", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := NewAuthenticatedSubject("firebase-uid-123", tt.authTime)
			if err != nil {
				t.Fatalf("NewAuthenticatedSubject: %v", err)
			}
			err = s.RequireRecentAuthentication(now)
			if tt.wantErr != (err != nil) {
				t.Fatalf("RequireRecentAuthentication() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrReauthenticationRequired) {
				t.Errorf("error = %v, want ErrReauthenticationRequired", err)
			}
		})
	}
}
