'use client'
import { useEffect } from "react";
import axios from "@/lib/axios/public";
import useSWR from "swr";
import { useAtom } from "jotai";
import { authAtom } from "@/atoms/auth";

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useAuth = () => {
    const [isAuthenticated, setIsAuthenticated] = useAtom(authAtom);
    const { data, isLoading, error } = useSWR(
        isAuthenticated ? '/api/users/me': null,
        fetcher
    );

    useEffect(() => {
        if (data) {
            setIsAuthenticated(true);
        } else if (error) {
            setIsAuthenticated(false);
        }
    }, [data, error, setIsAuthenticated]);

    return {
        isAuthenticated,
        user: data,
        isLoading,
        error
    }
};