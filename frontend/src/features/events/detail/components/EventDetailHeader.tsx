import React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import StatusBadge from '@/components/StatusBadge';
import { EVENT_STATUS_COLORS, EVENT_STATUS_LABELS } from '@/features/events/status';
import type { EventDraftDetail } from '@/features/events/types';
import { Pencil } from 'lucide-react';
import DeleteButton from './DeleteButton';

interface EventDetailHeaderProps {
    eventID: string;
    detail: EventDraftDetail;
}

const EventDetailHeader: React.FC<EventDetailHeaderProps> = ({ eventID, detail }) => {
    return (
        <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="min-w-0 space-y-2">
                <h1 className="break-words text-2xl font-bold leading-snug tracking-normal text-gray-900">
                    {detail.title}
                </h1>
                <StatusBadge
                    label={EVENT_STATUS_LABELS[detail.status]}
                    circleColor={EVENT_STATUS_COLORS[detail.status]}
                    textColor={EVENT_STATUS_COLORS[detail.status]}
                    textSize="sm"
                />
            </div>
            <div className="flex shrink-0 items-center gap-2">
                <Button variant="outline" asChild>
                    <Link href={`/events/${eventID}/edit`}>
                        <Pencil className="size-4" />
                        編集
                    </Link>
                </Button>
                <DeleteButton eventID={eventID} title={detail.title} />
            </div>
        </div>
    );
};

export default EventDetailHeader;
