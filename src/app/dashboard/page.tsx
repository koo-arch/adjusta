import React from 'react';
import Calendar from '@/features/calendar/Calendar';
import UpcomingEvents from '@/features/dashboard/UpcomingEvents';
import NeedsActionDrafts from '@/features/dashboard/NeedsActionDrafts';

const DashboardPage = () => {
    return (
        <div className="mx-auto max-w-screen-md p-4">
            <main className="grid grid-rows-[auto,auto] grid-cols-3 gap-4">
                <section className="col-span-3 grid grid-cols-2 gap-8">
                    <NeedsActionDrafts />
                    <UpcomingEvents />
                </section>
                <section className="col-span-3">
                    <h2 className="text-lg font-bold">カレンダー</h2>
                    <Calendar />
                </section>
            </main>
        </div>
    )
}

export default DashboardPage;