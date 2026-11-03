export function logRequestError(scope: string, error: unknown) {
  console.error(`[${scope}]`, error)
}
