'use client'
import type { ReactNode } from 'react';
import { QueryClientProvider, QueryClient, QueryCache, MutationCache, environmentManager } from "@tanstack/react-query";
import { getDefaultStore } from 'jotai/vanilla';
import { authErrorAtom } from '@/features/auth/store/error';
import { isUnauthorizedAPIError, isGoogleReauthorizationRequiredError } from '@/lib/api/errors';
import ToastProvider from './ToastProvider';

interface ProvidersProps {
    children: ReactNode;
}

// フルリロードでリセットされる。複数クエリの同時 401 でも redirect は一回だけにする
let hasRedirectedToLogin = false;

const handleGlobalAuthError = (error: unknown) => {
    // サーバー側 QueryClient には window がない
    if (typeof window === 'undefined') {
        return;
    }

    if (isUnauthorizedAPIError(error)) {
        if (hasRedirectedToLogin) {
            return;
        }

        hasRedirectedToLogin = true;
        // 期限切れ cookie を失効させてから /login に着地する
        // (直接 /login へ行くと cookie 保持中は proxy が /dashboard へ弾く)。
        // フルリロードになるため query cache / Jotai atom の認証済みデータも破棄される
        window.location.assign('/api/auth/session-expired');
        return;
    }

    if (isGoogleReauthorizationRequiredError(error)) {
        // Google 再認可要求は Adjusta のログイン失効とは別概念なので redirect しない
        getDefaultStore().set(authErrorAtom, {
            isOpen: true,
            message: 'Googleアカウントの再認可が必要です。再度ログインしてください。',
        });
    }
};

const makeQueryClient = () =>
    new QueryClient({
        // mutation 時のセッション切れも取りこぼさないよう QueryCache / MutationCache の両方に配線する
        queryCache: new QueryCache({ onError: handleGlobalAuthError }),
        mutationCache: new MutationCache({ onError: handleGlobalAuthError }),
        defaultOptions: {
            queries: {
                staleTime: 30_000,
                refetchOnWindowFocus: false,
                // 401 はリトライしても成功しないため即エラー確定し、redirect までの遅延を削る
                retry: (failureCount, error) =>
                    !isUnauthorizedAPIError(error) && failureCount < 1,
            },
        },
    });

let browserQueryClient: QueryClient | undefined;

const getQueryClient = () => {
    if (environmentManager.isServer()) {
        // サーバーはリクエストごとに新規（共有しない）
        return makeQueryClient();
    }

    // ブラウザはシングルトン（タブ内で使い回す）
    if (!browserQueryClient) {
        browserQueryClient = makeQueryClient();
    }

    return browserQueryClient;
};

const Providers = ({ children }: ProvidersProps) => {
    const queryClient = getQueryClient();

    return (
        <QueryClientProvider client={queryClient}>
            <ToastProvider>
                {children}
            </ToastProvider>
        </QueryClientProvider>
    )
}

export default Providers;
