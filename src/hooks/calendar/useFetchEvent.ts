'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtom } from 'jotai';
import type { GoogleEvent } from './type';

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useFetchEvent = () => {
    const [isAuthenticated] = useAtom(authAtom);
    const { data, isLoading, error } = useSWR<GoogleEvent[]>(
        isAuthenticated ? '/api/calendar/list' : null,
        fetcher
    );

    return {
        events: data,
        isLoading,
        error
    }
}

