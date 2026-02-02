import { defineConfig } from "orval";

export default defineConfig({
  backend: {
    input: {
      // путь из apps/frontend до контракта
      target: "../../libs/contracts/openapi/backend.yaml",
    },
    output: {
      mode: "split",
      target: "src/shared/api/generated",
      client: "fetch",
      clean: true,
    },
  },
});
