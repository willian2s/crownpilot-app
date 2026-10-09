package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Variáveis do Firebase Admin decididas na ADR 005. Os nomes são exatamente
// os do hosting (Render), sem prefixo CROWNPILOT_.
const (
	EnvFirebaseProjectID    = "FIREBASE_ADMIN_PROJECT_ID"
	EnvFirebaseClientEmail  = "FIREBASE_ADMIN_CLIENT_EMAIL"
	EnvFirebasePrivateKey   = "FIREBASE_ADMIN_PRIVATE_KEY"
	EnvFirebaseEmulatorHost = "FIREBASE_AUTH_EMULATOR_HOST"
)

// Firebase guarda a configuração do Firebase Admin. PrivateKey chega como veio
// do ambiente, com "\n" possivelmente escapado; a normalização acontece só na
// montagem do JSON em memória, no adapter.
type Firebase struct {
	ProjectID   string
	ClientEmail string
	PrivateKey  Secret

	// EmulatorHost espelha FIREBASE_AUTH_EMULATOR_HOST, que o SDK lê sozinho
	// do ambiente. Com ele definido o SDK aceita tokens SEM assinatura, por
	// isso a config o lê só para poder recusá-lo fora de Local.
	EmulatorHost string
}

// FirebaseMode é o modo de autenticação resultante da configuração.
type FirebaseMode int

const (
	// FirebaseDisabled: nenhuma variável definida. Só é válido em Local, onde
	// todo bearer é recusado com 401.
	FirebaseDisabled FirebaseMode = iota
	// FirebaseEmulator: Local apontando para o Firebase Auth Emulator.
	FirebaseEmulator
	// FirebaseCredentials: as três variáveis definidas; validação real.
	FirebaseCredentials
)

// Mode deriva o modo a partir dos campos. Só tem significado depois de
// Validate ter aceitado a configuração.
func (f Firebase) Mode() FirebaseMode {
	switch {
	case f.EmulatorHost != "":
		return FirebaseEmulator
	case f.ProjectID != "":
		return FirebaseCredentials
	default:
		return FirebaseDisabled
	}
}

func loadFirebase(lookup LookupFunc) Firebase {
	get := func(key string) string {
		v, _ := lookup(key)
		return v
	}
	return Firebase{
		ProjectID:    get(EnvFirebaseProjectID),
		ClientEmail:  get(EnvFirebaseClientEmail),
		PrivateKey:   NewSecret(get(EnvFirebasePrivateKey)),
		EmulatorHost: get(EnvFirebaseEmulatorHost),
	}
}

// projectIDPattern segue a regra de IDs de projeto do Google Cloud: 6 a 30
// caracteres, minúsculas, dígitos e hífen, começando por letra.
var projectIDPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{4,28}[a-z0-9]$`)

// validateFirebase aplica a política tudo-ou-nada da ADR 005. As mensagens
// citam só nomes de variáveis, nunca valores.
func validateFirebase(f Firebase, env Environment) error {
	var errs []error

	if f.ProjectID != "" && !projectIDPattern.MatchString(f.ProjectID) {
		errs = append(errs, fmt.Errorf("%s is not a valid Firebase project ID", EnvFirebaseProjectID))
	}

	if f.EmulatorHost != "" {
		if env != Local {
			// Fail-closed mais importante desta task: o emulator desliga a
			// verificação de assinatura dentro do SDK.
			return errors.Join(append(errs,
				fmt.Errorf("%s is allowed only in Local", EnvFirebaseEmulatorHost))...)
		}
		if f.ProjectID == "" {
			errs = append(errs, fmt.Errorf("%s is required with %s", EnvFirebaseProjectID, EnvFirebaseEmulatorHost))
		}
		if f.ClientEmail != "" || !f.PrivateKey.IsZero() {
			errs = append(errs, fmt.Errorf("%s and %s cannot be combined with %s",
				EnvFirebaseClientEmail, EnvFirebasePrivateKey, EnvFirebaseEmulatorHost))
		}
		return errors.Join(errs...)
	}

	set := 0
	for _, present := range []bool{f.ProjectID != "", f.ClientEmail != "", !f.PrivateKey.IsZero()} {
		if present {
			set++
		}
	}
	switch {
	case set == 0 && env != Local:
		errs = append(errs, fmt.Errorf("%s, %s and %s are required outside Local",
			EnvFirebaseProjectID, EnvFirebaseClientEmail, EnvFirebasePrivateKey))
	case set > 0 && set < 3:
		errs = append(errs, fmt.Errorf("%s, %s and %s must be configured together",
			EnvFirebaseProjectID, EnvFirebaseClientEmail, EnvFirebasePrivateKey))
	}

	if f.ClientEmail != "" && !strings.Contains(f.ClientEmail, "@") {
		errs = append(errs, fmt.Errorf("%s is not an email address", EnvFirebaseClientEmail))
	}

	return errors.Join(errs...)
}
