import React from 'react';
import { unstable_rethrow } from 'next/navigation';
import { User } from 'lucide-react';
import { requireUser } from '@/lib/server/api';
import UserButton from '@/features/auth/components/UserButton';
import type { AuthUser } from '@/features/auth/types';

// requireUser の redirect はストリーミング開始後でも機能する
// (クライアント側リダイレクトが注入される)ため、Suspense 内で await してよい
const UserMenu = async () => {
    let user: AuthUser | null = null;
    try {
        user = await requireUser();
    } catch (error) {
        // 401 → session-expired の NEXT_REDIRECT は Next に処理させる
        unstable_rethrow(error);
    }

    if (!user) {
        // layout 内のエラーは (app)/error.tsx に届かず既定のエラー画面で
        // ページ全体が落ちるため、backend 障害時は UserMenu だけ縮退させる
        return (
            <span
                aria-label="ユーザー情報を取得できませんでした"
                title="ユーザー情報を取得できませんでした"
                className="grid size-8 place-items-center rounded-full bg-muted"
            >
                <User className="size-4 text-muted-foreground" aria-hidden="true" />
            </span>
        );
    }

    return <UserButton user={user} />;
};

export default UserMenu;
