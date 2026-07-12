'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchNeedsActionDrafts } from '@/features/events/api/fetchNeedsActionDrafts';
import { buildNeedsActionDraftsQueryKey } from '@/features/events/queryKeys';

export const useFetchNeedsActionDrafts = () => {
    const { data, isPending, error, refetch } = useQuery({
        queryKey: buildNeedsActionDraftsQueryKey(),
        queryFn: fetchNeedsActionDrafts,
    });

    return {
        needsActionDrafts: data,
        isPending,
        error,
        refetch,
    };
}
