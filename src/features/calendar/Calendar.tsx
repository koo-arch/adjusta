'use client'
import React, { useEffect } from 'react';
import { useAtom } from 'jotai';
import { allEventsAtom } from '@/atoms/calendar';
import { StyleWrapper } from './style';
import FullCalendar from '@fullcalendar/react';
import type { ToolbarInput, DateRangeInput, EventClickArg, EventDropArg, DateSelectArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import timeGridPlugin from '@fullcalendar/timegrid';
import jaLocale from '@fullcalendar/core/locales/ja';
import momentPlugin from '@fullcalendar/moment';
import { useFetchGoogleEvent } from '@/hooks/calendar/useFetchGoogleEvent';
import { useSearchEvents } from '@/hooks/event/useSearchEvents';
import { renderDayCell, renderDayHeader, renderSlotLabel } from './render';
import type { CalendarEvent } from './type';
import type { EventDraftDetail } from '@/hooks/event/type';


interface CalendarProps<T extends CalendarEvent> {
    initialView?: "dayGridMonth" | "dayGridWeek" | "timeGridWeek" | "timeGridDay";
    firstDay?: number;
    headerToolbar?: ToolbarInput;
    select?: (arg: DateSelectArg) => void;
    selectedEvents?: T[];
    visibleRange?: DateRangeInput;
    eventClick?: (e: EventClickArg) => void;
    eventDrop?: (e: EventDropArg) => void;
    eventResize?: (e: EventResizeDoneArg) => void;
    editEvent?: EventDraftDetail
}

const Calendar = <T extends CalendarEvent>({
    initialView = 'dayGridMonth',
    firstDay,
    headerToolbar,
    select,
    selectedEvents,
    visibleRange,
    eventClick,
    eventDrop,
    eventResize,
    editEvent,
}: CalendarProps<T>) => {
    const { events, isLoading: isGoogleEventLoading, error: googleEventError } = useFetchGoogleEvent();
    const { searchEvents, isLoading: isSearchLoading, error: searchError } = useSearchEvents({ status: "pending" });
    const [allEvents, setAllEvents] = useAtom(allEventsAtom);


    useEffect(() => {
        if (!isGoogleEventLoading && !isSearchLoading) {
            const googleEventList: CalendarEvent[]  = events?.map(event => {
                return {
                    id: event.id,
                    title: event.summary,
                    start: event.start,
                    end: event.end,
                    location: event.location,
                    description: event.description,
                    origin: "google" as "google" | "local"
                };
            })
            ?.filter(event => {
                // 編集中のイベントは除外する
                if (editEvent) {
                    return !(editEvent.google_event_id === event.id);
                }
                return true;
            }) || [];
            
            const searchEventList: CalendarEvent[] = searchEvents?.flatMap(event => {
                return event.proposed_dates
                ?.filter(date => !editEvent?.proposed_dates?.some(edit => edit.id === date.id))
                .map(date => ({
                    id: date.id,
                    title: event.title,
                    start: date.start,
                    end: date.end,
                    location: event.location,
                    description: event.description,
                    origin: "local" as "google" | "local",
                    local_event_id: event.id
                }));
            }) || [];

            const allEvents: CalendarEvent[] = [...googleEventList, ...searchEventList, ...(selectedEvents || [])];
            setAllEvents(allEvents);
        }
    }, [events, searchEvents, isGoogleEventLoading, isSearchLoading, editEvent, setAllEvents, selectedEvents]);
    

    return (
        <div>
            <StyleWrapper>
                <FullCalendar
                    plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, momentPlugin]}
                    initialView={initialView}
                    firstDay={firstDay}
                    headerToolbar={headerToolbar || {
                        left: 'prev,next today',
                        center: 'title',
                        right: 'dayGridMonth,timeGridWeek,timeGridDay'
                    }}
                    businessHours={{ daysOfWeek: [1, 2, 3, 4, 5] }}
                    eventClick={eventClick || (() => {})}
                    snapDuration={'00:10:00'}
                    height={'auto'}
                    selectable={true}
                    selectMirror={true}
                    editable={true} // イベントのドラッグ＆ドロップを可能に
                    eventDrop={eventDrop || (() => {})}
                    eventResizableFromStart={true} // イベントの開始時間もリサイズ可能にする
                    eventResize={eventResize || (() => {})}
                    select={select || (() => {})}
                    visibleRange={visibleRange}
                    dayCellContent={renderDayCell}
                    dayHeaderContent={renderDayHeader}
                    slotLabelContent={renderSlotLabel}
                    aspectRatio={1.6}
                    locale={jaLocale}
                    events={allEvents}
                />
            </StyleWrapper>
        </div>
    );
}

export default Calendar;