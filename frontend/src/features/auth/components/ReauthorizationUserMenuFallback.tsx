'use client'

import { useEffect } from 'react';
import { User } from 'lucide-react';
import { useSetAtom } from 'jotai';
import { authErrorAtom } from '@/features/auth/store/error';

const ReauthorizationUserMenuFallback = () => {
    const setAuthError = useSetAtom(authErrorAtom);

    useEffect(() => {
        setAuthError({
            isOpen: true,
            message: 'Googleアカウントの再認可が必要です。',
        });
    }, [setAuthError]);

    return (
        <span
            aria-label="Googleアカウントの再認可が必要です"
            title="Googleアカウントの再認可が必要です"
            className="grid size-8 place-items-center rounded-full bg-muted"
        >
            <User className="size-4 text-muted-foreground" aria-hidden="true" />
        </span>
    );
};

export default ReauthorizationUserMenuFallback;
