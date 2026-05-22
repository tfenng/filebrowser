export type PasswordUpdateError =
  | "missingPassword"
  | "passwordMismatch"
  | "missingCurrentPassword"
  | "missingUser";

type PasswordUpdateInput = {
  password: string;
  passwordConf: string;
  currentPassword: string;
  requiresCurrentPassword: boolean;
  hasUser: boolean;
};

export function getPasswordUpdateError({
  password,
  passwordConf,
  currentPassword,
  requiresCurrentPassword,
  hasUser,
}: PasswordUpdateInput): PasswordUpdateError | null {
  if (!hasUser) return "missingUser";
  if (password === "") return "missingPassword";
  if (password !== passwordConf) return "passwordMismatch";
  if (requiresCurrentPassword && currentPassword === "") {
    return "missingCurrentPassword";
  }

  return null;
}
