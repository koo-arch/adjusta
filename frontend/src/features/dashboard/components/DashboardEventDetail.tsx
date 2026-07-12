'use client'
import React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import type { CalendarEvent } from '@/features/calendar/types';
import { formatJaDateSpan } from '@/lib/date/format';
import { CalendarDays, MapPin } from 'lucide-react';

interface DashboardEventDetailProps {
    event: CalendarEvent;
}

// カレンダーでクリックしたイベントの詳細を詳細タブ内に表示する(モーダルは使わない)
const DashboardEventDetail: React.FC<DashboardEventDetailProps> = ({ event }) => {
    return (
        <div className="space-y-3">
            <h2 className="min-w-0 break-words text-lg font-bold leading-snug tracking-normal text-gray-900">
                {event.title}
            </h2>
            {event.description && (
                <p className="whitespace-pre-wrap break-words text-sm text-gray-700">{event.description}</p>
            )}
            <div className="space-y-1.5">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <CalendarDays className="size-4 shrink-0" />
                    {formatJaDateSpan(event.start, event.end)}
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <MapPin className="size-4 shrink-0" />
                    {event.location || '未設定'}
                </div>
            </div>
            {event.origin === 'local' ? (
                <Button variant="ghost" className="px-0 text-primary hover:text-primary-dark" asChild>
                    <Link href={`/events/${event.local_event_id}`}>詳細ページへ</Link>
                </Button>
            ) : (
                <p className="text-xs text-muted-foreground">Google カレンダーの予定</p>
            )}
        </div>
    );
};

export default DashboardEventDetail;
