import { useSetAtom } from 'jotai';
import { authAtom } from '@/atoms/auth';

export const useLogin = () => {
    const setIsAuthenticated = useSetAtom(authAtom);

    return () => {
        setIsAuthenticated(true);
        window.location.href = 'http://localhost:8080/auth/google/login';
    }
}