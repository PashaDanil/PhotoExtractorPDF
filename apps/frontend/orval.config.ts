import { defineConfig } from 'orval';

export default defineConfig({
  imgpdf: {
    input: '../../libs/contracts/openapi/imgpdf.yaml',
    output: {
        client: 'react-query',
        target: './src/shared/api/generated/imgpdf.ts',
        schemas: './src/shared/api/generated/model',
    },
  },
});