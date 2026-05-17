'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtomValue } from 'jotai';
import { UpcomingEvent } from './type';

const fetcher = (url: string) => axios.get(url).then((res) => res.data);

export const useFetchUpcomingEvents = () => {
    const isAuthenticated = useAtomValue(authAtom);
    const { data, isLoading, error } = useSWR<UpcomingEvent[]>(
        isAuthenticated ? '/api/event/confirmed/upcoming' : null,
        fetcher
    );
   
    return {
        upcomingEvents: data,
        isLoading,
        error
    };
}