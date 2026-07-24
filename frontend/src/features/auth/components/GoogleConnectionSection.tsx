'use client'
import React from 'react';
import { useSetAtom } from 'jotai';
import StatusBadge from '@/components/common/StatusBadge/StatusBadge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { GOOGLE_CONNECTION_STATUS } from '@/features/auth/connectionStatus';
import { useAccounts } from '@/features/auth/hooks/useAccounts';
import { useLogout } from '@/features/auth/hooks/useLogout';
import { authErrorAtom } from '@/features/auth/store/error';

const ConnectionStatus = () => {
    const { connectionState, refetch } = useAccounts();
    const setAuthError = useSetAtom(authErrorAtom);

    if (connectionState.kind === 'loading') {
        return (
            <div className="space-y-2">
                <Skeleton className="h-5 w-24" />
                <Skeleton className="h-4 w-64" />
            </div>
        );
    }

    if (connectionState.kind === 'reauthorization_required') {
        const status = GOOGLE_CONNECTION_STATUS[connectionState.kind];

        return (
            <div className="space-y-3">
                <StatusBadge label={status.label} color={status.color} textSize="sm" />
                <p className="text-sm text-muted-foreground">
                    Google アカウントへのアクセス許可が失効しています。再認可すると引き続き Google
                    カレンダーと連携できます。
                </p>
                <Button
                    onClick={() =>
                        // 409 時の AuthErrorModal と同じ再認可フローに接続する(導線を一元化。screen-design 5.8)
                        setAuthError({
                            isOpen: true,
                            message: 'Googleアカウントの再認可が必要です。',
                        })
                    }
                >
                    再認可する
                </Button>
            </div>
        );
    }

    if (connectionState.kind === 'load_failed') {
        return (
            <div className="space-y-3">
                <p className="text-sm text-muted-foreground">連携状態を取得できませんでした。</p>
                <Button variant="outline" onClick={() => void refetch()}>
                    再試行
                </Button>
            </div>
        );
    }

    const status = GOOGLE_CONNECTION_STATUS[connectionState.kind];

    return (
        <div className="space-y-2">
            <div className="flex items-center gap-2">
                <StatusBadge label={status.label} color={status.color} textSize="sm" />
                <span className="text-sm text-muted-foreground">
                    {connectionState.account.email}
                </span>
            </div>
            <p className="text-sm text-muted-foreground">
                Google カレンダーの予定の取得と、確定した予定・候補日程の登録に利用しています。
            </p>
        </div>
    );
};

const GoogleConnectionSection = () => {
    const { logout, isPending } = useLogout();

    return (
        <Card>
            <CardHeader>
                <CardTitle>Google 連携</CardTitle>
                <CardDescription>
                    Adjusta は Google アカウントでログインし、Google カレンダーと連携します。
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <ConnectionStatus />
                <div className="border-t pt-4">
                    <AlertDialog>
                        <AlertDialogTrigger asChild>
                            <Button variant="outline" disabled={isPending}>
                                ログアウト
                            </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                            <AlertDialogHeader>
                                <AlertDialogTitle>ログアウトしますか?</AlertDialogTitle>
                                <AlertDialogDescription>
                                    ログアウトすると、再度 Google アカウントでのログインが必要になります。
                                </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                                <AlertDialogCancel>キャンセル</AlertDialogCancel>
                                <AlertDialogAction onClick={() => void logout()}>
                                    ログアウト
                                </AlertDialogAction>
                            </AlertDialogFooter>
                        </AlertDialogContent>
                    </AlertDialog>
                </div>
            </CardContent>
        </Card>
    );
};

export default GoogleConnectionSection;
