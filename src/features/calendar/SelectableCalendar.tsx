'use client'
import React, { useState, useRef } from 'react';
import { useAtom, useAtomValue } from 'jotai';
import type { EventClickArg, EventDropArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import Calendar from './Calendar';
import { EventImpl } from '@fullcalendar/core/internal';
import PopupMenu from '@/components/PopupMenu';
import type { CalendarEvent } from './type';
import type { EventDraftDetail } from '@/hooks/event/type';

type DateSelectInfo = {
    id: string;
    start: Date;
    end: Date;
}

interface SelectableCalendarProps<TDate extends DateSelectInfo, TEvent extends CalendarEvent> {
    dateAtom: any;
    eventAtom: any;
    editingEvent?: EventDraftDetail;
}

const SelectableCalendar = <TDate extends DateSelectInfo, TEvent extends CalendarEvent>({
    dateAtom,
    eventAtom,
    editingEvent,
}:SelectableCalendarProps<TDate, TEvent>) => {
    const [selectedDates, setSelectedDates] = useAtom<TDate[]>(dateAtom);
    const selectedEvents = useAtomValue<TEvent[]>(eventAtom);
    const [clickedEvent, setClickedEvent] = useState<EventImpl>();
    const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
    const buttonRef = useRef<HTMLButtonElement>(null);

    const handleDateSelect = (selectInfo: { start: Date, end: Date }) => {
        const newDate = {
            id : new Date().getTime().toString(),
            start: selectInfo.start,
            end: selectInfo.end,
        } as TDate;
        setSelectedDates([...selectedDates, newDate]);
    }

    // イベントをクリックした時にポップアップを表示する
    const handleEventClick = (e: EventClickArg) => {
        if (buttonRef.current) {
            buttonRef.current.click();
        }
        setClickedEvent(e.event);
        setPopupPosition({ top: e.jsEvent.clientY, left: e.jsEvent.clientX });
    }

    // イベントのドラッグ＆ドロップ時の処理
    const handleEventDrop = (e: EventDropArg) => {
        const updatedDates = selectedDates.map((date) => {
            if (date.id === e.event.id) {
                return {
                    ...date,
                    start: e.event.start || date.start,
                    end: e.event.end || date.end
                };
            }
            return date;
        });

        setSelectedDates(updatedDates);
    }

    // イベントの開始・終了時間の変更時の処理
    const handleEventResize = (e: EventResizeDoneArg) => {
        const updatedDates = selectedDates.map((date) => {
            if (date.id === e.event.id) {
                return {
                    ...date,
                    start: e.event.start || date.start,
                    end: e.event.end || date.end
                };
            }
            return date;
        });

        setSelectedDates(updatedDates);
    }

    // イベントの削除時の処理
    const handleDeleteEvent = () => {
        if (clickedEvent) {
            clickedEvent.remove();
            setSelectedDates((prev) => prev.filter((date) => date.id !== clickedEvent.id));
        }
    }

    return (
        <div>
            <Calendar 
                initialView="timeGridWeek"
                headerToolbar={{
                    left: 'prev,next today',
                    center: 'title',
                    right: '',
                }}
                select={handleDateSelect}
                selectedEvents={selectedEvents}
                eventClick={handleEventClick}
                eventDrop={handleEventDrop}
                eventResize={handleEventResize}
                editEvent={editingEvent}
            />
            <PopupMenu
                items={[
                    { label: '削除', onClick: handleDeleteEvent },
                ]}
                position={popupPosition}
                buttonRef={buttonRef}
            />
        </div>
    )
}

export default SelectableCalendar;