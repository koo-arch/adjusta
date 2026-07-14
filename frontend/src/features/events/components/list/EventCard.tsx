import React from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import StatusBadge from '@/components/common/StatusBadge/StatusBadge';
import { formatJaDateSpan } from '@/lib/date/format';
import { EVENT_STATUS_COLORS, EVENT_STATUS_LABELS } from '@/features/events/status';
import type { EventDraftDetail } from '@/features/events/types';
import { Calendar, ChevronRight } from 'lucide-react';

interface EventCardProps {
    event: EventDraftDetail;
}

const EventCard: React.FC<EventCardProps> = ({ event }) => {
    const confirmedDate = event.proposed_dates?.find((date) => date.id === event.confirmed_date_id);
    const isConfirmed = event.status === 'confirmed' && !!confirmedDate;
    const proposedDates = event.proposed_dates ?? [];

    return (
        <Link
            href={`/events/${event.id}`}
            className="block h-full rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
        >
            <Card className="h-full transition-shadow hover:shadow-md">
                <CardHeader className="flex-row items-start justify-between gap-2 space-y-0">
                    <CardTitle className="truncate">{event.title}</CardTitle>
                    <div className="shrink-0">
                        <StatusBadge
                            label={EVENT_STATUS_LABELS[event.status]}
                            color={EVENT_STATUS_COLORS[event.status]}
                            textSize="sm"
                        />
                    </div>
                </CardHeader>
                <CardContent>
                    {isConfirmed ? (
                        <div className="flex items-center text-sm text-muted-foreground">
                            <Calendar className="mr-1 size-4 shrink-0" />
                            {formatJaDateSpan(confirmedDate?.start, confirmedDate?.end)}
                        </div>
                    ) : proposedDates.length > 0 ? (
                        <div className="space-y-1">
                            {proposedDates.slice(0, 2).map((date) => (
                                <div key={date.id} className="flex items-center text-sm text-muted-foreground">
                                    <ChevronRight className="mr-1 size-4 shrink-0 text-yellow-500" />
                                    {formatJaDateSpan(date.start, date.end)}
                                </div>
                            ))}
                            {proposedDates.length > 2 && (
                                <p className="text-sm text-gray-400">... 他 {proposedDates.length - 2} 件</p>
                            )}
                        </div>
                    ) : (
                        <p className="text-sm text-gray-400">候補日程はありません</p>
                    )}
                </CardContent>
            </Card>
        </Link>
    );
};

export default EventCard;
