import React from 'react';
import Link from 'next/link';
import EventEdit from '@/features/events/edit/components/EventEdit';
import { ChevronLeft } from 'lucide-react';

interface EventEditPageContainerProps {
    eventID: string;
}

const EventEditPageContainer: React.FC<EventEditPageContainerProps> = ({ eventID }) => {
    return (
        <main className="mx-auto max-w-screen-lg space-y-4 px-4 py-8">
            <Link
                href={`/events/${eventID}`}
                className="inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
                <ChevronLeft className="size-4" />
                詳細へ戻る
            </Link>
            <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">イベント編集</h1>
            <EventEdit eventID={eventID} />
        </main>
    );
};

export default EventEditPageContainer;
