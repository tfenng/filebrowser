import { describe, expect, it } from "vitest";

import { shouldShowMetadataMismatch } from "../metadata";

describe("shouldShowMetadataMismatch", () => {
  it("shows metadata only for mismatched file types", () => {
    expect(
      shouldShowMetadataMismatch({
        extensionType: "image/jpeg",
        detectedType: "image/webp",
        typeMismatch: true,
      })
    ).toBe(true);

    expect(
      shouldShowMetadataMismatch({
        extensionType: "image/jpeg",
        detectedType: "image/jpeg",
        typeMismatch: false,
      })
    ).toBe(false);
  });
});
