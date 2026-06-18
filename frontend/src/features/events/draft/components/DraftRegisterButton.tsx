'use client'
import React from 'react';
import IconButton from '@/components/IconButton';
import { useAuth } from '@/hooks/auth/useAuth';
import { PlusIcon } from '@heroicons/react/20/solid';

const DraftRegisterButton = () => {
    const { isAuthenticated, user, isLoading } = useAuth();

    if (isLoading) return null;
    if (!isAuthenticated || !user) return null;

    return (
        <IconButton
            iconColor="primary"
            iconSize="lg"
            to="/schedule/draft/register"
        >
            <PlusIcon />
        </IconButton>
    )
}

export default DraftRegisterButton;