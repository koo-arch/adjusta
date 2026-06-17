'use client'
import { useAtom } from 'jotai';
import { fetchEventDetailAtomFamily } from '@/atoms/queries/event';


export const useFetchEventDetail = (eventID: string) => {
    const [{data, isLoading, error}] = useAtom(fetchEventDetailAtomFamily(eventID));

    return {
        eventDetail: data,
        isLoading,
        error
    };
}
