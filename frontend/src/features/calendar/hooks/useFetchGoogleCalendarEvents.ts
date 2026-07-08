'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchGoogleCalendarEvents } from '@/features/calendar/api/fetchGoogleCalendarEvents';
import { buildGoogleCalendarEventsQueryKey } from '@/features/calendar/queryKeys';

export const useFetchGoogleCalendarEvents = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildGoogleCalendarEventsQueryKey(),
        queryFn: fetchGoogleCalendarEvents,
    });

    return {
        events: data,
        isLoading,
        error,
    };
};
