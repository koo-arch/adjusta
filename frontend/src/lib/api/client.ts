type QueryValue = string | number | boolean | Date | null | undefined;

type QueryParams = object;

interface RequestOptions extends Omit<RequestInit, 'body'> {
    body?: BodyInit | object | null;
    query?: QueryParams;
    allowStatuses?: number[];
}

export interface APIClientResponse<T> {
    data: T;
    status: number;
    headers: Headers;
}

interface ParsedAPIErrorResponse {
    code?: string;
    error?: string;
}

export class APIClientError extends Error {
    constructor(
        message: string,
        readonly status: number,
        readonly data: unknown,
        readonly code?: string,
    ) {
        super(message);
        this.name = 'APIClientError';
    }
}

const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, '') ?? '';

const normalizePath = (path: string) => (path.startsWith('/') ? path : `/${path}`);

const buildURL = (path: string, query?: QueryParams) => {
    const url = `${baseURL}${normalizePath(path)}`;
    if (!query) {
        return url;
    }

    const searchParams = new URLSearchParams();

    for (const [key, rawValue] of Object.entries(query as Record<string, unknown>)) {
        const values = Array.isArray(rawValue) ? rawValue : [rawValue];

        for (const value of values) {
            if (
                value == null ||
                (
                    typeof value !== 'string' &&
                    typeof value !== 'number' &&
                    typeof value !== 'boolean' &&
                    !(value instanceof Date)
                )
            ) {
                continue;
            }

            searchParams.append(
                key,
                value instanceof Date ? value.toISOString() : String(value)
            );
        }
    }

    const search = searchParams.toString();
    return search ? `${url}?${search}` : url;
};

const buildBody = (body?: RequestOptions['body']) => {
    if (body == null) {
        return undefined;
    }

    if (
        body instanceof FormData ||
        body instanceof URLSearchParams ||
        body instanceof Blob ||
        typeof body === 'string'
    ) {
        return body;
    }

    return JSON.stringify(body);
};

const parseResponse = async <T>(response: Response): Promise<T> => {
    if (response.status === 204) {
        return undefined as T;
    }

    const contentType = response.headers.get('content-type') || '';
    if (contentType.includes('application/json')) {
        // ゲートウェイ等が application/json で空ボディを返すケースに耐える
        const text = await response.text();
        return (text ? JSON.parse(text) : undefined) as T;
    }

    return response.text() as Promise<T>;
};

const parseAPIErrorResponse = (data: unknown): ParsedAPIErrorResponse | undefined => {
    if (typeof data !== 'object' || data === null) {
        return undefined;
    }

    const code = 'code' in data && typeof data.code === 'string' ? data.code : undefined;
    const error = 'error' in data && typeof data.error === 'string' ? data.error : undefined;

    return code || error ? { code, error } : undefined;
};

const request = async <T>(path: string, options: RequestOptions = {}): Promise<APIClientResponse<T>> => {
    const { body, query, headers, allowStatuses = [], ...init } = options;

    const response = await fetch(buildURL(path, query), {
        ...init,
        credentials: 'include',
        headers: {
            ...(body instanceof FormData ? {} : { 'Content-Type': 'application/json' }),
            ...headers,
        },
        body: buildBody(body),
    });

    const data = await parseResponse<unknown>(response);
    const isAllowedStatus = allowStatuses.includes(response.status);

    if (!response.ok && !isAllowedStatus) {
        const apiError = parseAPIErrorResponse(data);
        const message = apiError?.error ?? `Request failed with status ${response.status}`;

        throw new APIClientError(message, response.status, data, apiError?.code);
    }

    return {
        data: data as T,
        status: response.status,
        headers: response.headers,
    };
};

export const apiClient = {
    get: <T>(path: string, options?: Omit<RequestOptions, 'method' | 'body'>) =>
        request<T>(path, { ...options, method: 'GET' }),
    post: <TResponse, TBody = unknown>(
        path: string,
        body?: TBody,
        options?: Omit<RequestOptions, 'method' | 'body'>,
    ) =>
        request<TResponse>(path, {
            ...options,
            method: 'POST',
            body: body as RequestOptions['body'],
        }),
    put: <TResponse, TBody = unknown>(
        path: string,
        body?: TBody,
        options?: Omit<RequestOptions, 'method' | 'body'>,
    ) =>
        request<TResponse>(path, {
            ...options,
            method: 'PUT',
            body: body as RequestOptions['body'],
        }),
    patch: <TResponse, TBody = unknown>(
        path: string,
        body?: TBody,
        options?: Omit<RequestOptions, 'method' | 'body'>,
    ) =>
        request<TResponse>(path, {
            ...options,
            method: 'PATCH',
            body: body as RequestOptions['body'],
        }),
    delete: <T>(path: string, options?: Omit<RequestOptions, 'method' | 'body'>) =>
        request<T>(path, { ...options, method: 'DELETE' }),
};
