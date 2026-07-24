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
import { authErrorAtom } from '@/features/auth/store/error';

const AuthErrorModal = () => {
    const [{ isOpen, message }, setAuthError] = useAtom(authErrorAtom);

    const handleClose = () => {
        setAuthError({ isOpen: false, message: '' });
    };

    const handleReauthorize = () => {
        const returnTo = `${window.location.pathname}${window.location.search}`;
        setAuthError({ isOpen: false, message: '' });
        window.location.assign(`/api/auth/google/reauthorize?return_to=${encodeURIComponent(returnTo)}`);
    };

    return (
        <Dialog open={isOpen} onOpenChange={(open) => !open && handleClose()}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Google連携の再認可</DialogTitle>
                    <DialogDescription>{message}</DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <Button variant="outline" onClick={handleClose}>あとで</Button>
                    <Button onClick={handleReauthorize}>Googleを再認可する</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

export default AuthErrorModal;
