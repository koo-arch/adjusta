'use client'
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '@/hooks/auth/useAuth';
import { fetchNeedsActionDrafts } from '@/features/events/api/fetchNeedsActionDrafts';
import { buildNeedsActionDraftsQueryKey } from '@/features/events/queryKeys';

export const useFetchNeedsActionDrafts = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useQuery({
        queryKey: buildNeedsActionDraftsQueryKey(),
        queryFn: fetchNeedsActionDrafts,
        enabled: isAuthenticated,
    });

    return {
        needsActionDrafts: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
}
