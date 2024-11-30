import React from 'react';
import ModalCalendar from '@/features/calendar/ModalCalendar';
import UpcomingEvents from '@/features/dashboard/UpcomingEvents';
import NeedsActionDrafts from '@/features/dashboard/NeedsActionDrafts';

const DashboardPage = () => {
    return (
        <div className="mx-auto max-w-screen-lg p-4">
            <main className="">
                <section className="mb-4">
                    <NeedsActionDrafts />
                </section>
                <section className="mb-4">
                    <UpcomingEvents />
                </section>
                <section className="mb-4">
                    <h2 className="text-lg font-bold">カレンダー</h2>
                    <ModalCalendar />
                </section>
            </main>
        </div>
    )
}

export default DashboardPage;