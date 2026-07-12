'use client'
import React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import type { CalendarEvent } from '@/features/calendar/types';
import { formatJaDateSpan } from '@/lib/date/format';
import { CalendarDays, MapPin, X } from 'lucide-react';

interface DashboardEventDetailProps {
    event: CalendarEvent;
    onClose: () => void;
}

// カレンダーでクリックしたイベントの詳細をパネル内に表示する(モーダルは使わない)
const DashboardEventDetail: React.FC<DashboardEventDetailProps> = ({ event, onClose }) => {
    return (
        <div className="space-y-3">
            <div className="flex items-start justify-between gap-2">
                <h2 className="min-w-0 break-words text-lg font-bold leading-snug tracking-normal text-gray-900">
                    {event.title}
                </h2>
                <Button
                    variant="ghost"
                    size="icon"
                    aria-label="詳細を閉じる"
                    title="詳細を閉じる"
                    className="h-8 w-8 shrink-0 text-muted-foreground hover:text-foreground"
                    onClick={onClose}
                >
                    <X />
                </Button>
            </div>
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
