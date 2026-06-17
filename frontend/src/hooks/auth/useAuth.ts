'use client'
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api/client";

export interface AuthUser {
    sub: string;
    name: string;
    email: string;
    picture: string;
}

export const currentUserKey = ['currentUser'] as const;

const fetchCurrentUser = async (): Promise<AuthUser | null> => {
    const response = await apiClient.get<AuthUser>('/api/users/me', {
        allowStatuses: [401],
    });

    if (response.status === 401) {
        return null;
    }

    return response.data ?? null;
};

export const useAuth = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: currentUserKey,
        queryFn: fetchCurrentUser,
    });

    return {
        isAuthenticated: !!data,
        user: data ?? null,
        isLoading,
        error
    }
};
