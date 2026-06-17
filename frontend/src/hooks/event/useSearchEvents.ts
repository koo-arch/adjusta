'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import type { EventDraftDetail, SearchParams } from './type';
import { useAuth } from '../auth/useAuth';

const fetchSearchEvents = async (params: SearchParams): Promise<EventDraftDetail[]> => {
    const response = await apiClient.get<EventDraftDetail[]>('/api/event/draft/search', {
        query: params,
    });
    return response.data;
};

export const useSearchEvents = (params: SearchParams) => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();

    const { data, isLoading, error } = useQuery({
        queryKey: ['draftEventSearch', params],
        queryFn: () => fetchSearchEvents(params),
        enabled: isAuthenticated,
    });

    return {
        searchEvents: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
};
