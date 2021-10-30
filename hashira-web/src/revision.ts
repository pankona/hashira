export const revision = (): string => {
  return process.env.REVISION || "unknown revision";
};
