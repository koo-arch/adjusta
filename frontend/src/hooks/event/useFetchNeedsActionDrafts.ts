'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { NeedsActionDraft } from './type';
import { useAuth } from '../auth/useAuth';

const fetcher = (url: string) => axios.get(url).then((res) => res.data);

export const useFetchNeedsActionDrafts = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useSWR<NeedsActionDraft[]>(
        isAuthenticated ? '/api/event/draft/needs-action' : null,
        fetcher
    );

    return {
        needsActionDrafts: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
}
