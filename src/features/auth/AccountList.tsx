'use client'
import React from 'react';
import { useAccounts } from '@/hooks/auth/useAccounts';

const AccountList: React.FC = () => {
    const { accounts, isLoading } = useAccounts();

    if (isLoading) return <div>Loading...</div>

    const accountList = accounts?.map(account => (
            <div key={account.account_id}>
                <p>{account.user_info.name}</p>
                <p>{account.user_info.email}</p>
            </div>
        ))

    return (
       <div>
            {accountList}
       </div>
    )
}

export default AccountList;