'use client'
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { fetchGoogleCalendarEvents } from '@/features/calendar/api/fetchGoogleCalendarEvents';
import { buildGoogleCalendarEventsQueryKey } from '@/features/calendar/queryKeys';

export const useFetchGoogleCalendarEvents = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: buildGoogleCalendarEventsQueryKey(),
        queryFn: fetchGoogleCalendarEvents,
        enabled: isAuthenticated,
    });

    return {
        events: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error,
    };
};
