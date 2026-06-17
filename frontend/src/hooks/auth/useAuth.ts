'use client'
import axios from "@/lib/axios/public";
import useSWR from "swr";

export interface AuthUser {
    sub: string;
    name: string;
    email: string;
    picture: string;
}

export const currentUserKey = '/api/users/me';

const fetcher = async (url: string): Promise<AuthUser | null> => {
    const response = await axios.get<AuthUser>(url, {
        validateStatus: (status) => status === 200 || status === 401,
    });

    if (response.status === 401) {
        return null;
    }

    return response.data;
};

export const useAuth = () => {
    const { data, isLoading, error } = useSWR<AuthUser | null>(
        currentUserKey,
        fetcher
    );

    return {
        isAuthenticated: !!data,
        user: data ?? null,
        isLoading,
        error
    }
};
