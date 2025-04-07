'use client'
import { useAtom } from 'jotai';
import { fetchEventDetailAtomFamily } from '@/atoms/queries/event';


export const useFetchEventDetail = (slug: string) => {
    const [{data, isLoading, error}] = useAtom(fetchEventDetailAtomFamily(slug));

    return {
        eventDetail: data,
        isLoading,
        error
    };
}