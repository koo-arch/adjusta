'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchAccount } from '@/features/auth/api/fetchAccount';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { buildAccountQueryKey } from '@/features/auth/queryKeys';

export const useAccounts = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: buildAccountQueryKey(),
        queryFn: fetchAccount,
        enabled: isAuthenticated,
    });

    return {
        account: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error,
    };
};
