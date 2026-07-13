'use client'

export const useLogin = () => {
    return () => {
        window.location.assign('/api/auth/google/login');
    };
};
