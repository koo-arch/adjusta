'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { AuthUser, useAuth } from './useAuth';

const fetchAccount = async () => {
    const response = await apiClient.get<AuthUser>('/api/account/list');
    return response.data;
};

export const useAccounts = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: ['account'],
        queryFn: fetchAccount,
        enabled: isAuthenticated,
    });

    return {
        account: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    }
}
