import { describe, expect, it } from "vitest";

import { normalizeLoginRedirect } from "../navigation";

describe("normalizeLoginRedirect", () => {
  it("uses the files view for missing, empty, or root redirects", () => {
    expect(normalizeLoginRedirect(undefined)).toBe("/files/");
    expect(normalizeLoginRedirect("")).toBe("/files/");
    expect(normalizeLoginRedirect("/")).toBe("/files/");
  });

  it("preserves explicit non-root redirects", () => {
    expect(normalizeLoginRedirect("/files/documents")).toBe(
      "/files/documents"
    );
    expect(normalizeLoginRedirect("/settings/profile")).toBe(
      "/settings/profile"
    );
  });
});
