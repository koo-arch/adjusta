import nextVitals from "eslint-config-next/core-web-vitals";
import storybook from "eslint-plugin-storybook";
import tanstackQuery from "@tanstack/eslint-plugin-query";

const config = [
  { ignores: ["storybook-static/", ".next-e2e/"] },
  ...nextVitals,
  ...storybook.configs["flat/recommended"],
  ...tanstackQuery.configs["flat/recommended"],
  {
    // Playwright の fixture は `use` コールバックを取るが React hook ではない
    files: ["e2e/**"],
    rules: { "react-hooks/rules-of-hooks": "off" },
  },
];

export default config;
