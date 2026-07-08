'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchDraftEventList } from '@/features/events/api/fetchDraftEventList';
import { buildDraftEventListQueryKey } from '@/features/events/queryKeys';

export const useFetchEventList = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildDraftEventListQueryKey(),
        queryFn: fetchDraftEventList,
    });

    return {
        eventList: data,
        isLoading,
        error
    };
};
