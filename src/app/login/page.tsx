import React from 'react';
import LoginButton from '@/features/auth/LoginButton';
import scheduleManage from '../../../public/images/schedule_manage.jpg';
import Card from '@/components/Card';

const LoignPage = () => {
    return (
        <div className="mx-auto max-w-screen-sm p-4 mb-4">
            <Card
                variant="outlined"
                image={scheduleManage}
            >
                <h1 className="text-center text-3xl font-extrabold mb-4">
                    Adjusta
                </h1>
                <p className="text-gray-500">
                    イベントの日程調整を簡単に行うことができる<wbr />サービスです。
                </p>
                <div className="mt-4 flex items-center justify-center">
                    <LoginButton />
                </div>
            </Card>
        </div>
    )
}

export default LoignPage;