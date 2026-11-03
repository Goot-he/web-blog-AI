export const env = {
  appName: import.meta.env.VITE_APP_NAME ?? 'CCEBD',
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? '/api',
  requestTimeout: Number(import.meta.env.VITE_REQUEST_TIMEOUT ?? 10000),
} as const
