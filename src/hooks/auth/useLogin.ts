import { useSetAtom } from 'jotai';
import { authAtom } from '@/atoms/auth';

export const useLogin = () => {
    const setIsAuthenticated = useSetAtom(authAtom);
    const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL;

    return () => {
        setIsAuthenticated(true);
        window.location.href = `${baseURL}/api/auth/google/login`;
    }
}