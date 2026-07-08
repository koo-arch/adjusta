import React from 'react';
import { requireUser } from '@/lib/server/api';
import UserButton from '@/features/auth/components/UserButton';

// requireUser の redirect はストリーミング開始後でも機能する
// (クライアント側リダイレクトが注入される)ため、Suspense 内で await してよい
const UserMenu = async () => {
    const user = await requireUser();

    return <UserButton user={user} />;
};

export default UserMenu;
