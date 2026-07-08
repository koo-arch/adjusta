'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchUpcomingEvents } from '@/features/events/api/fetchUpcomingEvents';
import { buildUpcomingEventsQueryKey } from '@/features/events/queryKeys';

export const useFetchUpcomingEvents = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildUpcomingEventsQueryKey(),
        queryFn: fetchUpcomingEvents,
    });

    return {
        upcomingEvents: data,
        isLoading,
        error
    };
}
