import React from 'react';
import Image from 'next/image';
import LoginButton from '@/features/auth/components/LoginButton';
import { Card, CardContent } from '@/components/ui/card';

const scheduleManageImage = '/images/schedule_manage.jpg';

const LoginPageContainer = () => {
    return (
        <div className="mx-auto max-w-screen-sm p-4 mb-4">
            <Card className="overflow-hidden">
                <div className="relative aspect-[16/9] w-full">
                    <Image
                        src={scheduleManageImage}
                        alt="スケジュール管理のイメージ"
                        fill
                        sizes="(max-width: 640px) 100vw, 640px"
                        className="object-cover"
                        priority
                    />
                </div>
                <CardContent className="p-6">
                    <h1 className="mb-4 text-center text-3xl font-extrabold">Adjusta</h1>
                    <p className="text-muted-foreground">
                        イベントの日程調整を簡単に行うことができる<wbr />サービスです。
                    </p>
                    <div className="mt-4 flex items-center justify-center">
                        <LoginButton />
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default LoginPageContainer;
