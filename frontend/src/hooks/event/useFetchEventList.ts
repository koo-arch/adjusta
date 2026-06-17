'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { EventDraftDetail } from './type'
import { useAuth } from '../auth/useAuth';

const fetcher = async (url: string) => axios.get<EventDraftDetail[]>(url).then((res) => res.data);

export const useFetchEventList = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useSWR<EventDraftDetail[]>(
        isAuthenticated ? '/api/calendar/event/draft/list' : null,
        fetcher
    );

    return {
        eventList: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
};
