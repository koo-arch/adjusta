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
import type { SelectedEvent } from '@/atoms/calendar';
import { renderDayCell, renderDayHeader, renderSlotLabel } from './render';

interface CalendarProps {
    initialView?: "dayGridMonth" | "timeGridWeek" | "timeGridDay";
    headerToolbar?: ToolbarInput;
    select?: (arg: { start: Date, end: Date }) => void;
    selectedEvents?: SelectedEvent[];
    eventClick?: (e: EventClickArg) => void;
    eventDrop?: (e: EventDropArg) => void;
    eventResize?: (e: EventResizeDoneArg) => void;
}

const Calendar: React.FC<CalendarProps> = ({ initialView, headerToolbar, select, selectedEvents, eventClick, eventDrop, eventResize }) => {
    const { events, isLoading, error } = useFetchEvent();

    if (isLoading) return <div>Loading...</div>;

    const eventList = events
        ?.filter(account => account?.events)
        ?.flatMap(account => account.events.map(event => ({
                id: event.id,
                title: event.summary,
                start: event.start,
                end: event.end,
                origin: "google"
            }
        ))
    );

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