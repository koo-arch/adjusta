'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import type { EventDraftDetail } from './type';

export const buildEventDetailQueryKey = (eventID?: string) => ['eventDetail', eventID] as const;

const fetchEventDetail = async (eventID: string) => {
    const response = await apiClient.get<EventDraftDetail>(`/api/calendar/event/draft/${eventID}`);
    return response.data;
};

export const useFetchEventDetail = (eventID: string) => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildEventDetailQueryKey(eventID),
        queryFn: () => fetchEventDetail(eventID),
        enabled: !!eventID,
    });

    return {
        eventDetail: data,
        isLoading,
        error
    };
}
