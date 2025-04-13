'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtom } from 'jotai';
import type {  GoogleCalendarResponse } from './type';

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useFetchGoogleEvent = () => {
    const [isAuthenticated] = useAtom(authAtom);
    const { data, isLoading, error } = useSWR<GoogleCalendarResponse>(
        isAuthenticated ? '/api/calendar/list' : null,
        fetcher
    );

    return {
        events: data,
        isLoading,
        error
    }
}

