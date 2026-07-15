'use client'
import React from 'react';
import { useAtom } from 'jotai';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { useLogout } from '@/features/auth/hooks/useLogout';
import { authErrorAtom } from '@/features/auth/store/error';

const AuthErrorModal = () => {
    const { logout } = useLogout();
    const [{ isOpen, message }, setAuthError] = useAtom(authErrorAtom);

    const handleClose = async () => {
        setAuthError({ isOpen: false, message: '' });
        // logout 成功時は useLogout 内の assign('/login') が遷移する(backend が cookie 破棄済み)。
        // 完了を await してから遷移することで、unload による logout リクエスト中断も避ける
        const ok = await logout();
        if (!ok) {
            // logout API 失敗時も cookie を失効させて /login に確実に着地させる
            window.location.assign('/api/auth/session-expired');
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={(open) => !open && void handleClose()}>
            <DialogContent onEscapeKeyDown={(event) => event.preventDefault()}>
                <DialogHeader>
                    <DialogTitle>認証エラー</DialogTitle>
                    <DialogDescription>{message}</DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <Button onClick={() => void handleClose()}>ログイン画面へ</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

export default AuthErrorModal;
