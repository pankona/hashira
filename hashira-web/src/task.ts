export const normalizeTasks = (lines: readonly string[]): string[] => {
  return lines.flatMap((line) => {
    const task = line.trim();
    return task ? [task] : [];
  });
};
