'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchAccount } from '@/features/auth/api/fetchAccount';
import { resolveGoogleConnectionState } from '@/features/auth/connectionStatus';
import { buildAccountQueryKey } from '@/features/auth/queryKeys';

export const useAccounts = () => {
    const { data, isLoading, error, refetch } = useQuery({
        queryKey: buildAccountQueryKey(),
        queryFn: fetchAccount,
    });
    const connectionState = resolveGoogleConnectionState({
        account: data,
        isLoading,
        error,
    });

    return {
        account: data,
        connectionState,
        isLoading,
        error,
        refetch,
    };
};
