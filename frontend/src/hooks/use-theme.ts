import { useEffect } from "react";

/**
 * Initialize dark mode on mount.
 * Checks localStorage or system preference.
 */
export function useInitTheme() {
  useEffect(() => {
    const stored = localStorage.getItem("theme");
    if (stored === "dark" || (!stored && window.matchMedia("(prefers-color-scheme: dark)").matches)) {
      document.documentElement.classList.add("dark");
    }
  }, []);
}
