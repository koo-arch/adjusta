'use client'
import React from 'react'
import { Button } from '@/components/ui/button';
import Image from 'next/image';
import { useLogin } from '@/features/auth/hooks/useLogin';

const LoginButton: React.FC = () => {
    const loginHandler = useLogin();

    return (
        <Button variant="outline" size="lg" className="rounded-full" onClick={loginHandler}>
            <Image src="/images/google.svg" alt="" height={24} width={24} />
            Googleでログイン
        </Button>
    )
}

export default LoginButton;
