package config

import (
	"strings"
	"testing"
)

func TestValidateFirebase(t *testing.T) {
	t.Parallel()

	full := testFirebase()
	without := func(mutate func(*Firebase)) Firebase {
		f := testFirebase()
		mutate(&f)
		return f
	}
	emulator := Firebase{ProjectID: "demo-crownpilot", EmulatorHost: "127.0.0.1:9099"}

	tests := []struct {
		name     string
		env      Environment
		firebase Firebase
		wantErr  []string // vazio = válido
	}{
		// Credenciais completas valem em qualquer ambiente.
		{name: "credentials in Local", env: Local, firebase: full},
		{name: "credentials in Preview", env: Preview, firebase: full},
		{name: "credentials in Staging", env: Staging, firebase: full},
		{name: "credentials in Production", env: Production, firebase: full},

		// Nada configurado: só Local, onde todo bearer vira 401.
		{name: "disabled in Local", env: Local, firebase: Firebase{}},
		{name: "disabled in Preview", env: Preview, firebase: Firebase{}, wantErr: []string{"required outside Local"}},
		{name: "disabled in Staging", env: Staging, firebase: Firebase{}, wantErr: []string{"required outside Local"}},
		{name: "disabled in Production", env: Production, firebase: Firebase{}, wantErr: []string{"required outside Local"}},

		// Tudo-ou-nada, inclusive em Local.
		{
			name: "missing private key", env: Production,
			firebase: without(func(f *Firebase) { f.PrivateKey = Secret{} }),
			wantErr:  []string{"must be configured together"},
		},
		{
			name: "missing client email", env: Local,
			firebase: without(func(f *Firebase) { f.ClientEmail = "" }),
			wantErr:  []string{"must be configured together"},
		},
		{
			name: "only project ID", env: Staging,
			firebase: Firebase{ProjectID: "crownpilot-test"},
			wantErr:  []string{"must be configured together"},
		},

		// Formatos.
		{
			name: "invalid project ID", env: Production,
			firebase: without(func(f *Firebase) { f.ProjectID = "Crown_Pilot" }),
			wantErr:  []string{EnvFirebaseProjectID + " is not a valid"},
		},
		{
			name: "invalid client email", env: Production,
			firebase: without(func(f *Firebase) { f.ClientEmail = "not-an-email" }),
			wantErr:  []string{EnvFirebaseClientEmail + " is not an email"},
		},

		// Emulator: só Local, com project ID e sem credencial real.
		{name: "emulator in Local", env: Local, firebase: emulator},
		{
			name: "emulator in Preview", env: Preview, firebase: emulator,
			wantErr: []string{EnvFirebaseEmulatorHost + " is allowed only in Local"},
		},
		{
			name: "emulator in Production even with credentials", env: Production,
			firebase: without(func(f *Firebase) { f.EmulatorHost = "127.0.0.1:9099" }),
			wantErr:  []string{EnvFirebaseEmulatorHost + " is allowed only in Local"},
		},
		{
			name: "emulator without project ID", env: Local,
			firebase: Firebase{EmulatorHost: "127.0.0.1:9099"},
			wantErr:  []string{EnvFirebaseProjectID + " is required with"},
		},
		{
			name: "emulator mixed with real credentials", env: Local,
			firebase: without(func(f *Firebase) { f.EmulatorHost = "127.0.0.1:9099" }),
			wantErr:  []string{"cannot be combined with " + EnvFirebaseEmulatorHost},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateFirebase(tt.firebase, tt.env)
			if len(tt.wantErr) == 0 {
				if err != nil {
					t.Fatalf("validateFirebase() error = %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("validateFirebase() = nil, want %q", tt.wantErr)
			}
			for _, want := range tt.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want it to contain %q", err, want)
				}
			}
		})
	}
}

func TestFirebaseErrorsNeverContainValues(t *testing.T) {
	t.Parallel()
	f := Firebase{
		ProjectID:   "crownpilot-test",
		ClientEmail: "leak-check@crownpilot-test.iam.gserviceaccount.com",
		PrivateKey:  NewSecret("PRIVATE-KEY-SENTINEL"),
	}
	f.EmulatorHost = "emulator-sentinel:9099"

	err := validateFirebase(f, Production)
	if err == nil {
		t.Fatal("validateFirebase() = nil, want error")
	}
	for _, value := range []string{"leak-check@", "PRIVATE-KEY-SENTINEL", "emulator-sentinel"} {
		if strings.Contains(err.Error(), value) {
			t.Errorf("error leaks a configured value %q: %v", value, err)
		}
	}
}

func TestFirebaseMode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		f    Firebase
		want FirebaseMode
	}{
		{"disabled", Firebase{}, FirebaseDisabled},
		{"credentials", testFirebase(), FirebaseCredentials},
		{"emulator", Firebase{ProjectID: "demo-crownpilot", EmulatorHost: "127.0.0.1:9099"}, FirebaseEmulator},
	}
	for _, tt := range tests {
		if got := tt.f.Mode(); got != tt.want {
			t.Errorf("%s: Mode() = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestLoadReadsFirebaseVariables(t *testing.T) {
	t.Parallel()
	cfg, err := Load(env(withFirebase(map[string]string{EnvEnvironment: "Production"})))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Firebase.Mode() != FirebaseCredentials {
		t.Errorf("Mode() = %d, want FirebaseCredentials", cfg.Firebase.Mode())
	}
	// A config guarda a chave como veio; normalizar "\n" é papel do adapter.
	if got := cfg.Firebase.PrivateKey.Reveal(); got != firebaseVars[EnvFirebasePrivateKey] {
		t.Errorf("PrivateKey changed during Load: %q", got)
	}
}
