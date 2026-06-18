'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { useAtom } from 'jotai';
import Modal from '@/components/Modal';
import { useLogout } from '@/features/auth/hooks/useLogout';
import { authErrorAtom } from '@/features/auth/store/error';

const AuthErrorModal = () => {
    const router = useRouter();
    const { logout } = useLogout();
    const [{ isOpen, message }, setAuthError] = useAtom(authErrorAtom);

    return (
        <Modal
            isOpen={isOpen}
            onClose={() => {
                setAuthError({ isOpen: false, message: '' });
                void logout();
                router.push('/login');
            }}
            title="認証エラー"
        >
            <p>{message}</p>
        </Modal>
    );
}

export default AuthErrorModal;
