'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import type { EventDraftDetail, SearchParams } from './type';
import { useAuth } from '../auth/useAuth';

const fetcher = async<T, U = undefined>(
    url: string,
    params: U,
): Promise<T> => {
    const { data } = await axios.get(url, { params });
    return data;
}

export const useSearchEvents = (params: SearchParams) => {
    const { isAuthenticated, isLoading: isAuthLoading, error: authError } = useAuth();

    const { data, isLoading, error } = useSWR<EventDraftDetail[]>(
        isAuthenticated? ['/api/event/draft/search',  params]: null,
        async ([url, params]) => fetcher(url, params)
    );

    return {
        searchEvents: data,
        isLoading: isAuthLoading || isLoading,
        error: authError ?? error
    };
};
