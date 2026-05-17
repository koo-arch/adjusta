'use client'
import React from 'react'
import Button from '@/components/Button';
import Image from 'next/image';
import { useLogin } from '@/hooks/auth/useLogin';

const LoginButton: React.FC = () => {
    const loginHandler = useLogin();

    return (
        <Button
            shape={'full'}
            variant='outline'
            intent='clear'
            size="lg"
            onClick={loginHandler}
            startIcon={
                <Image
                    src="https://www.svgrepo.com/show/475656/google-color.svg"
                    loading="lazy"
                    alt="google logo"
                    height={24}
                    width={24}
                ></Image>
            }
        >
            Googleでログイン
        </Button>
    )
}

export default LoginButton;