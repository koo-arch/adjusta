'use client'
import { useAtom } from 'jotai';
import { fetchEventDetailAtomFamily } from '@/atoms/queries/event';


export const useFetchEventDetail = (id: string) => {
    const [{data, isLoading, error}] = useAtom(fetchEventDetailAtomFamily(id));

    return {
        eventDetail: data,
        isLoading,
        error
    };
}