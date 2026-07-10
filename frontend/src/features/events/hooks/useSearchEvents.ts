'use client'
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import type { SearchParams } from '../types';
import { searchDraftEvents } from '@/features/events/api/searchDraftEvents';
import { buildDraftEventSearchQueryKey } from '@/features/events/queryKeys';

export const useSearchEvents = (params: SearchParams) => {
    const { data, isPending, isPlaceholderData, error, refetch } = useQuery({
        queryKey: buildDraftEventSearchQueryKey(params),
        queryFn: () => searchDraftEvents(params),
        // ページ・タブ切替時に前の結果を保持し、レイアウトの明滅を防ぐ
        placeholderData: keepPreviousData,
    });

    return {
        searchEvents: data?.items,
        pagination: data?.pagination,
        isPending,
        isPlaceholderData,
        error,
        refetch,
    };
};
