import React from 'react';
import EventCard from './EventCard';
import CreateEventPlaceholderCard from './CreateEventPlaceholderCard';
import type { EventDraftDetail } from '@/features/events/types';
import { cn } from '@/lib/utils';

interface EventGridProps {
    events: EventDraftDetail[];
    /** 疎な状態でグリッド末尾に破線の作成導線を出す */
    showCreatePlaceholder: boolean;
    /** 次ページ取得中は前ページの内容を薄めて表示する */
    isDimmed?: boolean;
}

const EventGrid: React.FC<EventGridProps> = ({ events, showCreatePlaceholder, isDimmed = false }) => (
    <ul
        className={cn(
            'grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3',
            isDimmed && 'opacity-60',
        )}
    >
        {events.map((event) => (
            <li key={event.id}>
                <EventCard event={event} />
            </li>
        ))}
        {showCreatePlaceholder && (
            <li className={cn(events.length > 0 && 'hidden lg:block')}>
                <CreateEventPlaceholderCard />
            </li>
        )}
    </ul>
);

export default EventGrid;
