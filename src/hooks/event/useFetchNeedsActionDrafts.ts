'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtomValue } from 'jotai';
import { NeedsActionDraft } from './type';

const fetcher = (url: string) => axios.get(url).then((res) => res.data);

export const useFetchNeedsActionDrafts = () => {
    const isAuthenticated = useAtomValue(authAtom);
    const { data, isLoading, error } = useSWR<NeedsActionDraft[]>(
        isAuthenticated ? '/api/event/draft/needs-action' : null,
        fetcher
    );

    return {
        needsActionDrafts: data,
        isLoading,
        error
    };
}