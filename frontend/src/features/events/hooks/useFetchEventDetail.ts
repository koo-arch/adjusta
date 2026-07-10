'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchEventDetail } from '@/features/events/api/fetchEventDetail';
import { buildEventDetailQueryKey } from '@/features/events/queryKeys';

export const useFetchEventDetail = (eventID: string) => {
    const { data, isPending, error, refetch } = useQuery({
        queryKey: buildEventDetailQueryKey(eventID),
        queryFn: () => fetchEventDetail(eventID),
        enabled: !!eventID,
    });

    return {
        eventDetail: data,
        isPending,
        error,
        refetch,
    };
}
