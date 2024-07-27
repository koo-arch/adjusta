import React from 'react';
import LoginButton from '@/features/auth/LoginButton';
import ThemeButton from '@/components/ThemeButton';

const LoignPage = () => {
    return (
        <div>
            <h1>Login</h1>
            <LoginButton />
            <ThemeButton />
        </div>
    )
}

export default LoignPage;