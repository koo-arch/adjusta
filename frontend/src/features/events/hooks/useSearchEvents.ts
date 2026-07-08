'use client'
import { useQuery } from '@tanstack/react-query';
import type { SearchParams } from '../types';
import { searchDraftEvents } from '@/features/events/api/searchDraftEvents';
import { buildDraftEventSearchQueryKey } from '@/features/events/queryKeys';

export const useSearchEvents = (params: SearchParams) => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildDraftEventSearchQueryKey(params),
        queryFn: () => searchDraftEvents(params),
    });

    return {
        searchEvents: data,
        isLoading,
        error
    };
};
