import { cache } from 'react';
import { headers } from 'next/headers';
import { redirect } from 'next/navigation';
import type { AuthUser } from '@/features/auth/types';

export class ServerAPIError extends Error {
    status: number;
    data: unknown;

    constructor(message: string, status: number, data: unknown) {
        super(message);
        this.name = 'ServerAPIError';
        this.status = status;
        this.data = data;
    }
}

// Server Component 専用の DAL。cookie を明示的に転送し、401 は cookie を
// 失効できる Route Handler に逃がす（RSC レンダリング中は Set-Cookie できない）。
export const serverApi = async <T>(path: string, init?: RequestInit): Promise<T> => {
    const baseURL = process.env.INTERNAL_BACKEND_URL?.replace(/\/$/, '');
    if (!baseURL) {
        throw new Error('INTERNAL_BACKEND_URL is not set');
    }

    // cookies().toString() は値を encodeURIComponent で再構築するため、
    // base64 パディング(=)を含む gin のセッション cookie が %3D に化けて
    // backend の検証が壊れる。ブラウザが送った Cookie ヘッダを生のまま転送する
    const headerStore = await headers();
    const response = await fetch(`${baseURL}${path}`, {
        ...init,
        headers: {
            cookie: headerStore.get('cookie') ?? '',
            ...init?.headers,
        },
        cache: 'no-store',
    });

    // redirect は NEXT_REDIRECT を throw するため、この関数を try/catch で包まないこと
    if (response.status === 401) {
        redirect('/api/auth/session-expired');
    }

    const text = await response.text();
    const data: unknown = text ? JSON.parse(text) : undefined;

    if (!response.ok) {
        const message =
            typeof data === 'object' &&
            data !== null &&
            'error' in data &&
            typeof data.error === 'string'
                ? data.error
                : `Request failed with status ${response.status}`;

        throw new ServerAPIError(message, response.status, data);
    }

    return data as T;
};

// React cache で同一リクエスト内の重複 /me 呼び出しをデデュープする
export const requireUser = cache(() => serverApi<AuthUser>('/api/users/me'));
