'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { AuthUser, useAuth } from './useAuth';

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useAccounts = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useSWR<AuthUser>(
        isAuthenticated ? '/api/account/list' : null,
        fetcher
    );

    return {
        account: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    }
}
