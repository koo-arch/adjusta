'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchCurrentUser } from '@/features/auth/api/fetchCurrentUser';
import { buildCurrentUserQueryKey } from '@/features/auth/queryKeys';

export const useAuth = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildCurrentUserQueryKey(),
        queryFn: fetchCurrentUser,
    });

    return {
        isAuthenticated: !!data,
        user: data ?? null,
        isLoading,
        error,
    };
};
