// orval.config.ts
import { defineConfig } from "orval"

export default defineConfig({
    api: {
        input: {
            target: "http://localhost:8080/swagger/doc.json",
            // если у тебя файл локально: target: "./openapi.yaml"
        },
        output: {
            mode: "split", // разнесёт по файлам (удобнее)
            target: "./src/shared/api/generated/generated.ts",
            schemas: "./src/shared/api/generated/model",
            client: "react-query",
            httpClient: "axios",
            clean: true,
            prettier: true,
            override: {
                // Генерим отдельный axios-инстанс (чтобы настроить baseURL, токены, interceptors)
                mutator: {
                    path: "./src/shared/api/http.ts",
                    name: "customInstance",
                },
            },
            tsconfig: "./tsconfig.app.json",
        },
    },
})