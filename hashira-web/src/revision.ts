export const revision = (): string => {
  return import.meta.env.VITE_REVISION || "unknown revision";
};
