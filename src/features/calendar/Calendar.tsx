'use client'
import React from 'react';
import { StyleWrapper } from './style';
import FullCalendar from '@fullcalendar/react';
import type { ToolbarInput, EventClickArg, EventDropArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import timeGridPlugin from '@fullcalendar/timegrid';
import jaLocale from '@fullcalendar/core/locales/ja';
import momentPlugin from '@fullcalendar/moment';
import { useFetchEvent } from '@/hooks/calendar/useFetchEvent';
import type { ProposedDate } from '@/atoms/calendar';
import { renderDayCell, renderDayHeader, renderSlotLabel } from './render';

type CalendarEvent = {
    id: string;
    title: string;
    start: Date;
    end: Date;
    origin: "google" | "local";
}

interface CalendarProps<T extends CalendarEvent> {
    initialView?: "dayGridMonth" | "timeGridWeek" | "timeGridDay";
    headerToolbar?: ToolbarInput;
    select?: (arg: { start: Date, end: Date }) => void;
    selectedEvents?: T[];
    eventClick?: (e: EventClickArg) => void;
    eventDrop?: (e: EventDropArg) => void;
    eventResize?: (e: EventResizeDoneArg) => void;
    editEvent?: ProposedDate[];
}

const Calendar = <T extends CalendarEvent>({
    initialView,
    headerToolbar,
    select,
    selectedEvents,
    eventClick,
    eventDrop,
    eventResize,
    editEvent,
}: CalendarProps<T>) => {
    const { events, isLoading, error } = useFetchEvent();

    if (isLoading) return <div>Loading...</div>;

    console.log(events);

    const eventList = events?.map(event => {
        return {
            id: event.id,
            title: event.summary,
            start: event.start,
            end: event.end,
            origin: "google"
        };
    })
        ?.filter(event => {
            // 編集中のイベントは除外する
            if (editEvent) {
                return !editEvent.some(date => date.event_id === event.id);
            }
            return true;
    })

    const conbinedEvents = selectedEvents && eventList ? [...eventList, ...selectedEvents] : eventList;

    return (
        <StyleWrapper>
            <FullCalendar
                plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, momentPlugin]}
                initialView={initialView || 'dayGridMonth'}
                headerToolbar={headerToolbar || {
                    left: 'prev,next today',
                    center: 'title',
                    right: 'dayGridMonth,timeGridWeek,timeGridDay'
                }}
                businessHours={{ daysOfWeek: [1, 2, 3, 4, 5] }}
                eventClick={eventClick || (() => {})}
                snapDuration={'00:10:00'}
                selectable={true}
                selectMirror={true}
                editable={true} // イベントのドラッグ＆ドロップを可能に
                eventDrop={eventDrop || (() => {})}
                eventResizableFromStart={true} // イベントの開始時間もリサイズ可能にする
                eventResize={eventResize || (() => {})}
                select={select || (() => {})}
                dayCellContent={renderDayCell}
                dayHeaderContent={renderDayHeader}
                slotLabelContent={renderSlotLabel}
                aspectRatio={1.6}
                locale={jaLocale}
                events={conbinedEvents}
            />
        </StyleWrapper>
    );
}

export default Calendar;