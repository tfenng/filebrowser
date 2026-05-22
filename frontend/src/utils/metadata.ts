export function shouldShowMetadataMismatch(
  metadata: FileMetadata | null | undefined
): boolean {
  return Boolean(
    metadata?.typeMismatch &&
    metadata.extensionType.trim() !== "" &&
    metadata.detectedType.trim() !== ""
  );
}
