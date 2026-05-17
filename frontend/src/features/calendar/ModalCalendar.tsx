'use client'
import React, { useState } from 'react';
import Link from 'next/link';
import { useAtomValue } from 'jotai';
import { allEventsAtom } from '@/atoms/calendar';
import Calendar from '@/features/calendar/Calendar';
import type { EventClickArg } from '@fullcalendar/core';
import Modal from '@/components/Modal/Modal';
import type { CalendarEvent } from '@/features/calendar/type';
import { CalendarDaysIcon, MapPinIcon } from "@heroicons/react/20/solid";
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

    const eventClose = () => {
        setIsModalOpen(false);
    }
    return (
        <div>
            <Calendar 
                eventClick={handleEventClick}
            />

            <Modal
                isOpen={isModalOpen}
                onClose={eventClose}
                title={selectedEvent?.title || ""}

            >
                <div>
                    {selectedEvent && (
                        <div>
                            <p className="text-sm text-gray-700 mt-4">{selectedEvent.description}</p>

                            <div className="flex item-center mt-2">
                                <CalendarDaysIcon className="h-5 w-5 mr-2" />
                                <p className="text-sm text-gray-500">
                                    {formatJaDateSpan(selectedEvent.start, selectedEvent.end)}
                                </p>

                            </div>
                            <div className="flex item-center mt-2">
                                <MapPinIcon className="h-5 w-5 mr-2" />
                                <p className="text-sm text-gray-500">
                                    {selectedEvent.location || "未設定"}
                                </p>
                            </div>
                            {selectedEvent.origin === "local" && (
                                <div className="mt-4">
                                    <Link 
                                        href={`/schedule/draft/${selectedEvent.slug}`}
                                        className="text-sm no-underline text-blue-500 hover:underline"
                                    >
                                       詳細ページ
                                    </Link>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </Modal>
        </div>
    )
}

export default ModalCalendar;