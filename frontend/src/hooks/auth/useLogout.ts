'use client'
import axios from '@/lib/axios/public';
import { mutate } from 'swr';
import { useRouter } from 'next/navigation';
import { currentUserKey } from './useAuth';

const fetcher = async (url: string) => await axios.get(url);

export const useLogout = () => {
    const router = useRouter();

    const logout = () => {
        fetcher('api/auth/logout')
            .then(() => {
                mutate(currentUserKey, null, false);
                mutate(
                    (key) => typeof key === 'string' && key.startsWith('/api/'),
                    undefined,
                    { revalidate: false }
                );
                router.push('/login');
            })
            .catch(err => console.error(err));
    }

    return logout 
}
