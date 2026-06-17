'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { UpcomingEvent } from './type';
import { useAuth } from '../auth/useAuth';

const fetcher = (url: string) => axios.get(url).then((res) => res.data);

export const useFetchUpcomingEvents = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useSWR<UpcomingEvent[]>(
        isAuthenticated ? '/api/event/confirmed/upcoming' : null,
        fetcher
    );
   
    return {
        upcomingEvents: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
}
