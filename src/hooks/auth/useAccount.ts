'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtom } from 'jotai';

interface Account {
    sub: string;
    name: string;
    email: string;
    picture: string;
}

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useAccounts = () => {
    const [isAuthenticated] = useAtom(authAtom);
    const { data, isLoading, error } = useSWR<Account>(
        isAuthenticated ? '/api/account/list' : null,
        fetcher
    );

    return {
        account: data,
        isLoading,
        error
    }
}