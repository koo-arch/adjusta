import React, { Suspense } from 'react';
import ProfileSection, { ProfileSectionSkeleton } from '@/features/auth/components/ProfileSection';
import GoogleConnectionSection from '@/features/auth/components/GoogleConnectionSection';
import CalendarSettingsSection from '@/features/auth/components/CalendarSettingsSection';
import { requireUser } from '@/lib/server/api';

const AccountProfile = async () => {
    const user = await requireUser();

    return <ProfileSection user={user} />;
};

const AccountPageContainer = () => {
    return (
        <main className="mx-auto max-w-screen-md space-y-6 px-4 py-8">
            <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">
                アカウント設定
            </h1>
            <Suspense fallback={<ProfileSectionSkeleton />}>
                <AccountProfile />
            </Suspense>
            <GoogleConnectionSection />
            <CalendarSettingsSection />
        </main>
    );
};

export default AccountPageContainer;
