'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchDraftEventList } from '@/features/events/api/fetchDraftEventList';
import { buildDraftEventListQueryKey } from '@/features/events/queryKeys';
import type { EventListParams } from '@/features/events/types';

export const useFetchEventList = (params: EventListParams = {}) => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildDraftEventListQueryKey(params),
        queryFn: () => fetchDraftEventList(params),
    });

    return {
        eventList: data?.items,
        pagination: data?.pagination,
        isLoading,
        error
    };
};
