'use client'
import React from 'react';
import { useAuth } from '@/hooks/auth/useAuth';

const UserInfo: React.FC = () => {
    const { user, isAuthenticated } = useAuth();
    console.log(isAuthenticated)

    if (!isAuthenticated) return <div>Not authenticated</div>

    return (
        <div>
            <p>{user?.picture}</p>
            <p>{user?.name}</p>
            <p>{user?.email}</p>
        </div>
    )
}

export default UserInfo;