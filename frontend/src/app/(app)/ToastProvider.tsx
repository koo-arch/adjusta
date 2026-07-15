'use client'
import React from 'react';
import { Toaster } from '@/components/ui/sonner';

interface ToastProviderProps {
    children: React.ReactNode;
}

const ToastProvider: React.FC<ToastProviderProps> = ({ children }) => {
    return (
        <>
            {children}
            <Toaster position="top-right" duration={5000} closeButton richColors />
        </>
    )
}

export default ToastProvider;
