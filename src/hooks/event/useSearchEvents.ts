'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtomValue } from 'jotai';
import type { EventDraftDetail, SearchParams } from './type';

const fetcher = async<T, U = undefined>(
    url: string,
    params: U,
): Promise<T> => {
    const { data } = await axios.get(url, { params });
    return data;
}

export const useSearchEvents = (params: SearchParams) => {
    const isAuthenticated = useAtomValue(authAtom);

    const { data, isLoading, error } = useSWR<EventDraftDetail[]>(
        isAuthenticated? ['/api/event/draft/search',  params]: null,
        async ([url, params]) => fetcher(url, params)
    );

    return {
        searchEvents: data,
        isLoading,
        error
    };
};