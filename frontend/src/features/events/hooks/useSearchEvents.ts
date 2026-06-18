'use client'
import { useQuery } from '@tanstack/react-query';
import type { SearchParams } from '../types';
import { useAuth } from '@/hooks/auth/useAuth';
import { searchDraftEvents } from '@/features/events/api/searchDraftEvents';
import { buildDraftEventSearchQueryKey } from '@/features/events/queryKeys';

export const useSearchEvents = (params: SearchParams) => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();

    const { data, isLoading, error } = useQuery({
        queryKey: buildDraftEventSearchQueryKey(params),
        queryFn: () => searchDraftEvents(params),
        enabled: isAuthenticated,
    });

    return {
        searchEvents: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
};
