'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { UpcomingEvent } from './type';
import { useAuth } from '../auth/useAuth';

const fetchUpcomingEvents = async () => {
    const response = await apiClient.get<UpcomingEvent[]>('/api/event/confirmed/upcoming');
    return response.data;
};

export const useFetchUpcomingEvents = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: ['upcomingEvents'],
        queryFn: fetchUpcomingEvents,
        enabled: isAuthenticated,
    });
   
    return {
        upcomingEvents: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
}
