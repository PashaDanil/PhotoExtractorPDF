import axios, { type AxiosRequestConfig } from "axios"

export const AXIOS_INSTANCE = axios.create({
  // В dev используется VITE_API_URL (http://backend:8080)
  // В production nginx проксирует /api/ на backend, поэтому используем относительный путь
  baseURL: import.meta.env.VITE_API_URL ?? "/api",
  withCredentials: false, // true если куки/сессии
})

// Пример: токен из localStorage
AXIOS_INSTANCE.interceptors.request.use((config) => {
  const token = localStorage.getItem("token")
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const customInstance = <T>(config: AxiosRequestConfig): Promise<T> => {
  return AXIOS_INSTANCE.request<T>(config).then((res) => res.data)
}
