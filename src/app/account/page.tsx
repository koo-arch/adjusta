import React from 'react';
import AddAccountButton from '@/features/auth/AddAccountButton';
import AccountList from '@/features/auth/AccountList';

const AccountPage = () => {
    return (
        <div>
            <h1>Account</h1>
            <AddAccountButton />
            <AccountList />
        </div>
    )
}

export default AccountPage;