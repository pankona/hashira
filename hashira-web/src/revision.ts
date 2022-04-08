export const revision = (): string => {
  return import.meta.env.REVISION || "unknown revision";
};
