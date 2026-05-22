import { describe, expect, it } from "vitest";

import { getPasswordUpdateError } from "../password";

describe("getPasswordUpdateError", () => {
  it("reports mismatched new passwords", () => {
    expect(
      getPasswordUpdateError({
        password: "new-password-123",
        passwordConf: "different-password-123",
        currentPassword: "current-password-123",
        requiresCurrentPassword: true,
        hasUser: true,
      })
    ).toBe("passwordMismatch");
  });

  it("reports a missing current password when required", () => {
    expect(
      getPasswordUpdateError({
        password: "new-password-123",
        passwordConf: "new-password-123",
        currentPassword: "",
        requiresCurrentPassword: true,
        hasUser: true,
      })
    ).toBe("missingCurrentPassword");
  });

  it("allows a complete password update request", () => {
    expect(
      getPasswordUpdateError({
        password: "new-password-123",
        passwordConf: "new-password-123",
        currentPassword: "current-password-123",
        requiresCurrentPassword: true,
        hasUser: true,
      })
    ).toBeNull();
  });
});
