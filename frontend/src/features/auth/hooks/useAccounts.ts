'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchAccount } from '@/features/auth/api/fetchAccount';
import { buildAccountQueryKey } from '@/features/auth/queryKeys';

export const useAccounts = () => {
    const { data, isLoading, error, refetch } = useQuery({
        queryKey: buildAccountQueryKey(),
        queryFn: fetchAccount,
    });

    return {
        account: data,
        isLoading,
        error,
        refetch,
    };
};
