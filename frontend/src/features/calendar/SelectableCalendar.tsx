'use client'
import React, { useState, useRef } from 'react';
import type { EventClickArg, EventDropArg, DateSelectArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import Calendar from './Calendar';
import { EventImpl } from '@fullcalendar/core/internal';
import PopupMenu from '@/components/PopupMenu';
import type { CalendarEvent } from './type';
import type { EventDraftDetail } from '@/features/events/types';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';

type DateSelectInfo = {
    id: string;
    start: Date;
    end: Date;
}

interface SelectableCalendarProps<TDate extends DateSelectInfo> {
    dates: TDate[];
    onDatesChange: React.Dispatch<React.SetStateAction<TDate[]>>;
    selectedEvents: CalendarEvent[];
    editingEvent?: EventDraftDetail;
}

const SelectableCalendar = <TDate extends SelectedDate | ProposedDate>({
    dates,
    onDatesChange,
    selectedEvents,
    editingEvent,
}:SelectableCalendarProps<TDate>) => {
    const [clickedEvent, setClickedEvent] = useState<EventImpl>();
    const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
    const buttonRef = useRef<HTMLButtonElement>(null);

    const handleDateSelect = (e: DateSelectArg) => {
        const newDate = {
            id : new Date().getTime().toString(),
            start: e.start,
            end: e.end,
        } as TDate;
        onDatesChange((prev) => [...prev, newDate]);
    }

    // イベントをクリックした時にポップアップを表示する
    const handleEventClick = (e: EventClickArg) => {
        const event = dates.find((date) => date.id === e.event.id);

        if (event) {
            if (buttonRef.current) {
                buttonRef.current.click();
            }

            setClickedEvent(e.event);
            setPopupPosition({ top: e.jsEvent.pageY, left: e.jsEvent.pageX });
        }
    }

    // イベントのドラッグ＆ドロップ時の処理
    const handleEventDrop = (e: EventDropArg) => {
        const updatedDates = dates.map((date) => {
            if (date.id === e.event.id) {
                return {
                    ...date,
                    start: e.event.start || date.start,
                    end: e.event.end || date.end
                };
            }
            return date;
        });

        onDatesChange(updatedDates);
    }

    // イベントの開始・終了時間の変更時の処理
    const handleEventResize = (e: EventResizeDoneArg) => {
        const updatedDates = dates.map((date) => {
            if (date.id === e.event.id) {
                return {
                    ...date,
                    start: e.event.start || date.start,
                    end: e.event.end || date.end
                };
            }
            return date;
        });

        onDatesChange(updatedDates);
    }

    // イベントの削除時の処理
    const handleDeleteEvent = () => {
        if (clickedEvent) {
            clickedEvent.remove();
            onDatesChange((prev) => prev.filter((date) => date.id !== clickedEvent.id));
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
