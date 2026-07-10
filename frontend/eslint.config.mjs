import nextVitals from "eslint-config-next/core-web-vitals";
import storybook from "eslint-plugin-storybook";
import tanstackQuery from "@tanstack/eslint-plugin-query";

const config = [
  { ignores: ["storybook-static/"] },
  ...nextVitals,
  ...storybook.configs["flat/recommended"],
  ...tanstackQuery.configs["flat/recommended"],
];

export default config;
