'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import type {  GoogleCalendarResponse } from './type';
import { useAuth } from '../auth/useAuth';

const fetchGoogleEvents = async () => {
    const response = await apiClient.get<GoogleCalendarResponse>('/api/calendar/list');
    return response.data;
};

export const useFetchGoogleEvent = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: ['googleCalendarEvents'],
        queryFn: fetchGoogleEvents,
        enabled: isAuthenticated,
    });

    return {
        events: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    }
}
