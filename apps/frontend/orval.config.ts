import { defineConfig } from "orval";

export default defineConfig({
  backend: {
    input: {
      target: "../../libs/contracts/openapi/backend.yaml",
    },
    output: {
      mode: "split",
      target: "extractor/shared/api/generated",
      client: "fetch",
      clean: true,
      override: {
        mutator: {
          path: "extractor/shared/api/mutator.ts",
          name: "customFetch",
        },
      },
    },
  },
});
