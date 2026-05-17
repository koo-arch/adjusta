'use client'
import React from 'react';
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { ReactQueryStreamedHydration } from '@tanstack/react-query-next-experimental';

interface ProvidersProps {
    children: React.ReactNode;
}

const queryClient = new QueryClient();

const Providers: React.FC<ProvidersProps> = ({ children }) => {
    return (
        <QueryClientProvider client={queryClient}>
            <ReactQueryStreamedHydration>
                {children}
            </ReactQueryStreamedHydration>
        </QueryClientProvider>
    )
}

export default Providers;