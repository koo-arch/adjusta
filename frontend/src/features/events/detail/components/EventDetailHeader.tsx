import React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import StatusBadge from '@/components/StatusBadge';
import { EVENT_STATUS_COLORS, EVENT_STATUS_LABELS } from '@/features/events/status';
import type { EventDraftDetail } from '@/features/events/types';
import { MapPin, Pencil } from 'lucide-react';
import DeleteButton from './DeleteButton';

interface EventDetailHeaderProps {
    eventID: string;
    detail: EventDraftDetail;
}

const EventDetailHeader: React.FC<EventDetailHeaderProps> = ({ eventID, detail }) => {
    return (
        <header className="space-y-3">
            <div className="flex items-start justify-between gap-4">
                <h1 className="min-w-0 break-words text-2xl font-bold leading-snug tracking-normal text-gray-900">
                    {detail.title}
                </h1>
                <div className="flex shrink-0 items-center gap-1">
                    <Button
                        variant="ghost"
                        size="icon"
                        className="text-muted-foreground hover:text-foreground [&_svg]:size-5"
                        asChild
                    >
                        <Link href={`/events/${eventID}/edit`} aria-label="編集" title="編集">
                            <Pencil />
                        </Link>
                    </Button>
                    <DeleteButton eventID={eventID} title={detail.title} />
                </div>
            </div>
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1">
                <StatusBadge
                    label={EVENT_STATUS_LABELS[detail.status]}
                    circleColor={EVENT_STATUS_COLORS[detail.status]}
                    textColor={EVENT_STATUS_COLORS[detail.status]}
                    textSize="sm"
                />
                {detail.location && (
                    <span className="flex items-center gap-1 text-sm text-muted-foreground">
                        <MapPin className="size-4 shrink-0" />
                        {detail.location}
                    </span>
                )}
            </div>
            {detail.description && (
                <p className="whitespace-pre-wrap break-words text-sm text-gray-700">{detail.description}</p>
            )}
        </header>
    );
};

export default EventDetailHeader;
