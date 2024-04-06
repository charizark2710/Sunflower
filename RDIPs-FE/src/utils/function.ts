export const initValue = (a: string | undefined | null): string => {
  if (a && Boolean(a)) return a;
  return '';
};
