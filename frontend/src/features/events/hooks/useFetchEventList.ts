'use client'
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '@/hooks/auth/useAuth';
import { fetchDraftEventList } from '@/features/events/api/fetchDraftEventList';
import { buildDraftEventListQueryKey } from '@/features/events/queryKeys';

export const useFetchEventList = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: buildDraftEventListQueryKey(),
        queryFn: fetchDraftEventList,
        enabled: isAuthenticated,
    });

    return {
        eventList: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
};
