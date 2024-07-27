'use client'
import axios from 'axios';
import { useRouter } from 'next/navigation';
import { authAtom } from '@/atoms/auth';
import { useSetAtom } from 'jotai';

const fetcher = async (url: string) => await axios.get(url);

export const useLogout = () => {
    const setIsAuthenticated = useSetAtom(authAtom);
    const router = useRouter();

    const logout = () => {
        fetcher('/api/auth/logout')
            .then(() => {
                setIsAuthenticated(false);
                router.push('/login');
            })
            .catch(err => console.error(err));
    }

    return logout 
}