'use client'
import axios from '@/lib/axios/public';
import { useRouter } from 'next/navigation';
import { authAtom } from '@/atoms/auth';
import { useSetAtom } from 'jotai';

const fetcher = async (url: string) => await axios.get(url);

export const useLogout = () => {
    const setIsAuthenticated = useSetAtom(authAtom);
    const router = useRouter();

    const logout = () => {
        fetcher('/auth/logout')
            .then(() => {
                setIsAuthenticated(false);
                router.push('/login');
            })
            .catch(err => console.error(err));
    }

    return logout 
}