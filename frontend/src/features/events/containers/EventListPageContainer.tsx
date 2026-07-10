import React, { Suspense } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import EventList, { EventListSkeleton } from '@/features/events/components/list/EventList';
import { Plus } from 'lucide-react';

const EventListPageContainer = () => {
    return (
        <main className="mx-auto max-w-screen-lg space-y-6 px-4 py-8">
            <div className="flex items-center justify-between gap-4">
                <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">
                    イベント一覧
                </h1>
                <Button asChild>
                    <Link href="/events/new">
                        <Plus className="size-4" />
                        イベントを作成
                    </Link>
                </Button>
            </div>
            {/* EventList は useSearchParams を使うため Suspense 境界が必要 */}
            <Suspense fallback={<EventListSkeleton />}>
                <EventList />
            </Suspense>
        </main>
    );
};

export default EventListPageContainer;
