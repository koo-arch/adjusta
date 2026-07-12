'use client'
import React, { useRef, useState } from 'react';
import { useAtomValue } from 'jotai';
import { allEventsAtom } from '@/features/events/store/calendar';
import Calendar from '@/features/calendar/components/Calendar';
import type { EventClickArg } from '@fullcalendar/core';
import { Dialog, DialogContent, DialogTitle } from '@/components/ui/dialog';
import { Popover, PopoverAnchor, PopoverContent } from '@/components/ui/popover';
import type { CalendarEvent } from '@/features/calendar/types';
import DashboardEventDetail from './DashboardEventDetail';

const DashboardCalendar: React.FC = () => {
    const allEvents = useAtomValue(allEventsAtom);
    const [selectedEvent, setSelectedEvent] = useState<CalendarEvent | null>(null);
    const [isPopoverOpen, setIsPopoverOpen] = useState(false);
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    // virtualRef は非 null の RefObject を要求する。popover を開く前に必ず代入される
    const anchorRef = useRef<HTMLElement>(null!);

    const handleEventClick = (e: EventClickArg) => {
        const event = allEvents.find(event => event.id === e.event.id);
        if (!event) {
            return;
        }
        setSelectedEvent(event);
        // クリック位置の近くに反応を出す: lg 以上はクリックしたイベント要素に
        // anchor した popover、それ未満は画面が狭いため dialog で表示する
        if (window.matchMedia('(min-width: 1024px)').matches) {
            anchorRef.current = e.el;
            setIsPopoverOpen(true);
        } else {
            setIsDialogOpen(true);
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

            <Popover open={isPopoverOpen} onOpenChange={setIsPopoverOpen}>
                <PopoverAnchor virtualRef={anchorRef} />
                <PopoverContent className="w-80" collisionPadding={8}>
                    {selectedEvent && <DashboardEventDetail event={selectedEvent} />}
                </PopoverContent>
            </Popover>

            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogContent className="max-w-md">
                    {/* Radix の a11y 要件(DialogTitle 必須)。見た目のタイトルは詳細側の h2 が担う */}
                    <DialogTitle className="sr-only">{selectedEvent?.title ?? 'イベント詳細'}</DialogTitle>
                    {selectedEvent && <DashboardEventDetail event={selectedEvent} />}
                </DialogContent>
            </Dialog>
        </div>
    );
}

export default DashboardCalendar;
