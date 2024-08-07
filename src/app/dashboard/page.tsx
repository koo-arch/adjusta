import React from 'react';
import UserInfo from '@/features/auth/UserInfo';
import Calendar from '@/features/calendar/Calendar';

const DashboardPage = () => {
    return (
        <div>
            <h1>Dashboard</h1>
            <UserInfo />
            <Calendar />
        </div>
    )
}

export default DashboardPage;