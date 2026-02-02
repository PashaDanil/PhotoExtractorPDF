export const customFetch = async <T>(
  url: string,
  options?: RequestInit
): Promise<T> => {
  const res = await fetch(url, options);

  // Для ZIP файлов возвращаем blob
  const contentType = res.headers.get('content-type');
  if (contentType?.includes('application/zip')) {
    return (await res.blob()) as T;
  }

  // Для пустых ответов
  if ([204, 205, 304].includes(res.status)) {
    return {} as T;
  }

  // Для JSON
  const text = await res.text();
  return text ? JSON.parse(text) : ({} as T);
};
