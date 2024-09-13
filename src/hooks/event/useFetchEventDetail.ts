'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { useAtomValue } from 'jotai';
import { authAtom } from '@/atoms/auth';
import { EventDraftDetail } from './type';

const fetcher = async (url: string) => axios.get<EventDraftDetail>(url).then((res) => res.data);

export const useFetchEventDetail = (id: string) => {
    const isAuthenticated = useAtomValue(authAtom);
    const { data, isLoading, error } = useSWR<EventDraftDetail>(
        isAuthenticated ? `/api/calendar/event/draft/${id}` : null,
        fetcher
    );

    return {
        eventDetail: data,
        isLoading,
        error
    };
}