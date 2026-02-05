const API_PREFIX = "/api";

export const customFetch = async <T>(
  url: string,
  options?: RequestInit
): Promise<T> => {
  const fullUrl = url.startsWith("http") ? url : `${API_PREFIX}${url}`;

  const res = await fetch(fullUrl, options);

  const contentType = res.headers.get("content-type") || "";

  // ZIP отдаём как Blob
  if (contentType.includes("application/zip")) {
    return (await res.blob()) as T;
  }

  if ([204, 205, 304].includes(res.status)) {
    return {} as T;
  }

  // JSON (или пусто)
  const text = await res.text();
  return text ? (JSON.parse(text) as T) : ({} as T);
};
