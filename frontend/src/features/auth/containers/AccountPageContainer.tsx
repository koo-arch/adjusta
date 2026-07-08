import React from 'react';
import UserInfo from '@/features/auth/components/UserInfo';
import { requireUser } from '@/lib/server/api';

const AccountPageContainer = async () => {
    const user = await requireUser();

    return (
        <div className="mx-auto max-w-screen-md p-4">
            <UserInfo user={user} />
        </div>
    );
};

export default AccountPageContainer;
