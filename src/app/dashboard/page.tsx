import React from 'react';
import Calendar from '@/features/calendar/Calendar';
import UpcomingEvents from '@/features/dashboard/UpcomingEvents';

const DashboardPage = () => {
    return (
        <div>
            <UpcomingEvents />
            <Calendar />
        </div>
    )
}

export default DashboardPage;