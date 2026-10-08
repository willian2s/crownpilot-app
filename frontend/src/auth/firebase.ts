import { getApp, getApps, initializeApp, type FirebaseApp, type FirebaseOptions } from 'firebase/app';
import { connectAuthEmulator, getAuth, type Auth } from 'firebase/auth';

const requiredConfigKeys = [
  'apiKey',
  'authDomain',
  'projectId',
  'appId',
] as const;

let firebaseAuth: Auth | null = null;
let emulatorConnected = false;

function isLoopbackHost(host: string): boolean {
  return host === 'localhost' || host === '127.0.0.1' || host === '[::1]' || host === '::1';
}

function isValidPort(port: string): boolean {
  return /^\d{1,5}$/.test(port) && Number(port) > 0 && Number(port) <= 65535;
}

function publicEnv(name: string): string {
  const value = import.meta.env[name];
  return typeof value === 'string' ? value.trim() : '';
}

export function getFirebaseConfig(): FirebaseOptions | null {
  const config: FirebaseOptions = {
    apiKey: publicEnv('VITE_FIREBASE_API_KEY'),
    authDomain: publicEnv('VITE_FIREBASE_AUTH_DOMAIN'),
    projectId: publicEnv('VITE_FIREBASE_PROJECT_ID'),
    storageBucket: publicEnv('VITE_FIREBASE_STORAGE_BUCKET') || undefined,
    messagingSenderId: publicEnv('VITE_FIREBASE_MESSAGING_SENDER_ID') || undefined,
    appId: publicEnv('VITE_FIREBASE_APP_ID'),
    measurementId: publicEnv('VITE_FIREBASE_MEASUREMENT_ID') || undefined,
  };

  return requiredConfigKeys.every((key) => Boolean(config[key])) ? config : null;
}

export function getFirebaseAuth(): Auth | null {
  if (firebaseAuth) {
    return firebaseAuth;
  }

  const config = getFirebaseConfig();
  if (!config) {
    return null;
  }

  const app: FirebaseApp = getApps().length > 0 ? getApp() : initializeApp(config);
  firebaseAuth = getAuth(app);

  const emulatorHost = publicEnv('VITE_FIREBASE_EMULATOR_HOST') || '127.0.0.1';
  const emulatorPort = publicEnv('VITE_FIREBASE_EMULATOR_PORT') || '9099';
  const emulatorRequested = publicEnv('VITE_FIREBASE_USE_EMULATOR').toLowerCase() === 'true';
  if (!emulatorConnected && import.meta.env.DEV && emulatorRequested &&
      isLoopbackHost(emulatorHost) && isValidPort(emulatorPort)) {
    const emulatorUrl = new URL(`http://${emulatorHost}:${emulatorPort}`);
    if (emulatorUrl.protocol === 'http:' && isLoopbackHost(emulatorUrl.hostname) &&
        !emulatorUrl.username && !emulatorUrl.password) {
      connectAuthEmulator(firebaseAuth, emulatorUrl.toString(), { disableWarnings: true });
    }
    emulatorConnected = true;
  }

  return firebaseAuth;
}
