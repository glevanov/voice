import { en } from "./locales/en";
import { sv } from "./locales/sv";

export type Lang = "en" | "sv";

/**
 * Generic type to get a path of a nested object
 */
type Path<Type> = Type extends object
  ? {
      [Key in keyof Type]: Key extends string
        ? Type[Key] extends object
          ? `${Key}` | `${Key}.${Path<Type[Key]>}`
          : `${Key}`
        : never;
    }[keyof Type]
  : never;

const traverse = (obj: Locale, path: Path<Locale>): string => {
  const parts = path.split(".");

  let value: Locale | Locale[keyof Locale] | string = obj;

  for (const part of parts) {
    if (typeof value !== "string") {
      // @ts-expect-error difficult to type
      value = value[part];
    }
  }

  return value as string;
};

export type Locale = typeof sv;

const locales: Record<Lang, Locale> = {
  en,
  sv,
};

const createI18n =
  (lang: Lang) =>
  (path: Path<Locale>): string => {
    return traverse(locales[lang], path);
  };

const language = import.meta.env.VITE_LANGUAGE as Lang;

export const i18n = createI18n(language);
