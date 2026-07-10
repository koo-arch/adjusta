import React from 'react';
import ProfileSection from '@/features/auth/components/ProfileSection';
import GoogleConnectionSection from '@/features/auth/components/GoogleConnectionSection';
import CalendarSettingsSection from '@/features/auth/components/CalendarSettingsSection';
import { requireUser } from '@/lib/server/api';

const AccountPageContainer = async () => {
    const user = await requireUser();

    return (
        <main className="mx-auto max-w-screen-md space-y-6 px-4 py-8">
            <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">
                アカウント設定
            </h1>
            <ProfileSection user={user} />
            <GoogleConnectionSection />
            <CalendarSettingsSection />
        </main>
    );
};

export default AccountPageContainer;
