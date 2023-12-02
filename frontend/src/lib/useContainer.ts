import { ref, watchEffect } from "vue";
import config from "@/lib/env";

type ApiError = { result: "ok" | "error"; message: string }

interface FetchResponse<T> {
  data: T | null;
  error: string | null;
  loading: boolean;
}

export async function useContainer<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<FetchResponse<T>> {
  const headers = new Headers();
  headers.set("Content-Type", "application/json");

  try {
    const response = await fetch(
      `http://localhost:${config.port}${endpoint}`,
      {
        ...options,
        headers,
      },
    );

    const responseData = await response.json();

    if (!response.ok) {
      throw new Error(responseData.message);
    }

    return {
      data: responseData as T,
      error: null,
      loading: false,
    };
  } catch (err: any) {
    return {
      data: null,
      error: err.message,
      loading: false,
    };
  }
}