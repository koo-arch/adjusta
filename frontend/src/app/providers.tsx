'use client'
import type { ReactNode } from 'react';
import { QueryClientProvider, QueryClient, isServer } from "@tanstack/react-query";
import ToastProvider from './ToastProvider';

interface ProvidersProps {
    children: ReactNode;
}

const makeQueryClient = () =>
    new QueryClient({
        defaultOptions: {
            queries: {
                staleTime: 30_000,
                refetchOnWindowFocus: false,
                retry: 1,
            },
        },
    });

let browserQueryClient: QueryClient | undefined;

const getQueryClient = () => {
    if (isServer) {
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
