'use client'
import useSWR from 'swr';
import axios from '@/lib/axios/public';
import { authAtom } from '@/atoms/auth';
import { useAtom } from 'jotai';

interface Accounts {
    account_id: string;
    user_info: AccountInfo;
}

interface AccountInfo {
    sub: string;
    name: string;
    email: string;
    picture: string;
}

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useAccounts = () => {
    const [isAuthenticated] = useAtom(authAtom);
    const { data, isLoading, error } = useSWR<Accounts[]>(
        isAuthenticated ? '/api/account/list' : null,
        fetcher
    );

    return {
        accounts: data,
        isLoading,
        error
    }
}