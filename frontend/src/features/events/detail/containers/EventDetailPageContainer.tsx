import React from 'react';
import Link from 'next/link';
import EventDetail from '@/features/events/detail/components/EventDetail';
import { ChevronLeft } from 'lucide-react';

interface EventDetailPageContainerProps {
    eventID: string;
}

const EventDetailPageContainer: React.FC<EventDetailPageContainerProps> = ({ eventID }) => {
    return (
        <main className="mx-auto max-w-screen-md space-y-4 px-4 py-8">
            <Link
                href="/events"
                className="inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
                <ChevronLeft className="size-4" />
                イベント一覧へ
            </Link>
            <EventDetail eventID={eventID} />
        </main>
    );
};

export default EventDetailPageContainer;
