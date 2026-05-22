const defaultFilesPath = "/files/";

export function normalizeLoginRedirect(redirect: unknown): string {
  if (typeof redirect !== "string") {
    return defaultFilesPath;
  }

  const trimmed = redirect.trim();
  if (trimmed === "" || trimmed === "/") {
    return defaultFilesPath;
  }

  return trimmed;
}
