'use client'
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { fetchUpcomingEvents } from '@/features/events/api/fetchUpcomingEvents';
import { buildUpcomingEventsQueryKey } from '@/features/events/queryKeys';

export const useFetchUpcomingEvents = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: buildUpcomingEventsQueryKey(),
        queryFn: fetchUpcomingEvents,
        enabled: isAuthenticated,
    });
   
    return {
        upcomingEvents: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
}
