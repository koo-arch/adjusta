'use client'
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api/client';
import { NeedsActionDraft } from './type';
import { useAuth } from '../auth/useAuth';

const fetchNeedsActionDrafts = async () => {
    const response = await apiClient.get<NeedsActionDraft[]>('/api/event/draft/needs-action');
    return response.data;
};

export const useFetchNeedsActionDrafts = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: ['needsActionDrafts'],
        queryFn: fetchNeedsActionDrafts,
        enabled: isAuthenticated,
    });

    return {
        needsActionDrafts: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
}
