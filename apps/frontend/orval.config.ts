import { defineConfig } from 'orval';

export default defineConfig({
  imgpdf: {
    input: '../../libs/contracts/openapi/imgpdf.yaml',
    output: {
        client: 'react-query',
        target: './src/api/imgpdf.ts',
        schemas: './src/api/model',
    },
  },
});