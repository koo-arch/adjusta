'use client'
import React, { useState } from 'react';
import Link from 'next/link';
import { useAtomValue } from 'jotai';
import { allEventsAtom } from '@/features/events/store/calendar';
import Calendar from '@/features/calendar/components/Calendar';
import type { EventClickArg } from '@fullcalendar/core';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import type { CalendarEvent } from '@/features/calendar/types';
import { CalendarDays, MapPin } from 'lucide-react';
import { formatJaDateSpan } from '@/lib/date/format';


const ModalCalendar: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedEvent, setSelectedEvent] = useState<CalendarEvent | null>(null);
    const allEvents = useAtomValue(allEventsAtom);

    const handleEventClick = (e: EventClickArg) => {
        const event = allEvents.find(event => event.id === e.event.id);
        if (event) {
            setSelectedEvent(event);
            setIsModalOpen(true);
        }
    }

    return (
        <div>
            {/* ダッシュボードのカレンダーは閲覧専用(選択・ドラッグ・リサイズ不可。ui-review P1 #3) */}
            <Calendar
                selectable={false}
                editable={false}
                eventClick={handleEventClick}
            />

            <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
                <DialogContent className="max-w-md">
                    <DialogHeader>
                        <DialogTitle className="break-words">{selectedEvent?.title || ''}</DialogTitle>
                    </DialogHeader>
                    {selectedEvent && (
                        <div className="space-y-2">
                            {selectedEvent.description && (
                                <p className="whitespace-pre-wrap break-words text-sm text-gray-700">
                                    {selectedEvent.description}
                                </p>
                            )}
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <CalendarDays className="size-4 shrink-0" />
                                {formatJaDateSpan(selectedEvent.start, selectedEvent.end)}
                            </div>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <MapPin className="size-4 shrink-0" />
                                {selectedEvent.location || '未設定'}
                            </div>
                        </div>
                    )}
                    {selectedEvent?.origin === 'local' && (
                        <DialogFooter>
                            <Button variant="ghost" className="text-primary hover:text-primary-dark" asChild>
                                <Link href={`/events/${selectedEvent.local_event_id}`}>詳細ページへ</Link>
                            </Button>
                        </DialogFooter>
                    )}
                </DialogContent>
            </Dialog>
        </div>
    )
}

export default ModalCalendar;
