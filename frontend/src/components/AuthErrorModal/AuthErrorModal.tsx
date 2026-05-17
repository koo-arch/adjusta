'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { useAtom } from 'jotai';
import { authErrorAtom } from '@/atoms/error';
import Modal from '@/components/Modal';
import { useLogout } from '@/hooks/auth/useLogout';

const AuthErrorModal = () => {
    const router = useRouter();
    const logout = useLogout();
    const [{ isOpen, message }, setAuthError] = useAtom(authErrorAtom);

    return (
        <Modal
            isOpen={isOpen}
            onClose={() => {
                setAuthError({ isOpen: false, message: '' });
                logout();
                router.push('/login');
            }}  
            title="認証エラー"
        >
            <p>{message}</p>
        </Modal>
    );
}

export default AuthErrorModal;