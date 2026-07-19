import React from 'react';
import { Skeleton } from '@/components/ui/skeleton';

// UserButton の trigger(h-8 w-8 の丸アイコン)と同寸に保ち、
// 静的シェル → ストリーミング差し替え時のレイアウトシフトを防ぐ
const UserMenuSkeleton: React.FC = () => {
    return (
        <Skeleton
            className="size-8 rounded-full"
            role="status"
            aria-label="ユーザー情報を読み込み中"
        />
    );
};

export default UserMenuSkeleton;
