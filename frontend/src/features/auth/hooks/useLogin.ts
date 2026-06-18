'use client'

export const useLogin = () => {
    const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL;

    return () => {
        window.location.href = `${baseURL}/api/auth/google/login`;
    };
};
