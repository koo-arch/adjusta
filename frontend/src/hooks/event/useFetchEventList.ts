'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { EventDraftDetail } from './type'
import { useAuth } from '../auth/useAuth';

const fetchEventList = async () => {
    const response = await apiClient.get<EventDraftDetail[]>('/api/calendar/event/draft/list');
    return response.data;
};

export const useFetchEventList = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: ['draftEventList'],
        queryFn: fetchEventList,
        enabled: isAuthenticated,
    });

    return {
        eventList: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
};
