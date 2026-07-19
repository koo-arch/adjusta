'use client'
import React, { startTransition } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';

// PPR では静的シェル送信後に動的ホール(requireUser 等)の例外が届くため、
// error boundary がないとページ全体が Next のデフォルトエラー画面に置き換わる
const AppError = ({
    reset,
}: {
    error: Error & { digest?: string };
    reset: () => void;
}) => {
    const router = useRouter();

    // reset() 単体はクライアント側の再レンダーのみで Server Component の
    // エラーから復帰できないため、refresh とあわせて実行する
    const retry = () => {
        startTransition(() => {
            router.refresh();
            reset();
        });
    };

    return (
        <main className="mx-auto max-w-screen-md px-4 py-8">
            <div className="flex flex-col items-center gap-4 py-16 text-center">
                <p className="text-sm text-muted-foreground">
                    ページの表示に失敗しました。時間をおいて再度お試しください。
                </p>
                <Button variant="outline" onClick={retry}>
                    再試行
                </Button>
            </div>
        </main>
    );
};

export default AppError;
