'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import type {  GoogleCalendarResponse } from './type';
import { useAuth } from '../auth/useAuth';

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useFetchGoogleEvent = () => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();
    const { data, isLoading, error } = useSWR<GoogleCalendarResponse>(
        isAuthenticated ? '/api/calendar/list' : null,
        fetcher
    );

    return {
        events: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    }
}
