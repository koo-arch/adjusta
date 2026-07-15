import React, { Suspense } from 'react';
import EventList, { EventListSkeleton } from '@/features/events/components/list/EventList';
import CreateEventFab from '@/features/events/components/list/CreateEventFab';

const EventListPageContainer = () => {
    return (
        // pb-24 は FAB とページネーション・最終行の重なりを防ぐ(md 以上は FAB なし)
        <main className="mx-auto max-w-6xl space-y-6 px-4 pb-24 pt-8 sm:px-6 md:pb-8">
            <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">
                イベント一覧
            </h1>
            {/* EventList は useSearchParams を使うため Suspense 境界が必要 */}
            <Suspense fallback={<EventListSkeleton />}>
                <EventList />
            </Suspense>
            <CreateEventFab />
        </main>
    );
};

export default EventListPageContainer;
