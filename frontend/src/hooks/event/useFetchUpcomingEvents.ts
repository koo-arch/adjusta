'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { UpcomingEvent } from './type';
import { useAuth } from '../auth/useAuth';

export const buildUpcomingEventsQueryKey = () => ['upcomingEvents'] as const;

const fetchUpcomingEvents = async () => {
    const response = await apiClient.get<UpcomingEvent[]>('/api/event/confirmed/upcoming');
    return response.data;
};

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
