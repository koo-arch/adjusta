'use client'
import React from 'react';
import Card from '@/components/Card';
import StatusBadge from '@/components/StatusBadge';
import { formatJaDateSpan } from '@/lib/date/format';
import { EventDraftDetail } from '@/hooks/event/type';
import { ChevronRightIcon } from '@heroicons/react/20/solid';
import { CalendarIcon } from '@heroicons/react/20/solid';
import { MapPinIcon } from '@heroicons/react/20/solid';

interface EventCardProps {
    event: EventDraftDetail;
    onClick: () => void;
}

const statusLabel = (status: string) => {
    switch (status) {
        case 'pending':
            return '調整中';
        case 'confirmed':
            return '確定';
        case 'rejected':
            return 'キャンセル';
        default:
            return '';
    }
}

const statusColor = (status: string) => {
    switch (status) {
        case 'pending':
            return 'yellow';
        case 'confirmed':
            return 'green';
        case 'rejected':
            return 'red';
        default:
            return 'gray';
    }
}

const EventCard: React.FC<EventCardProps> = ({ event, onClick }) => {
    const confirmedDate= event.proposed_dates?.find((date) => date.id === event.confirmed_date_id);
    const isConfirmed = event.status === 'confirmed' && !!event.confirmed_date_id && !!confirmedDate;

    return (
        <Card variant="outlined" background="inherit" isButton={true} onClick={onClick}>
            <div>
                <div className="flex justify-between items-center mb-2">
                    <h2 className="text-lg font-bold">{event?.title}</h2>
                    <StatusBadge 
                        label={statusLabel(event.status)}
                        circleColor={statusColor(event.status)}
                        textSize="sm"
                     />
                </div>
      
                {isConfirmed ? (
                    <div className="flex items-center text-sm text-gray-500 mb-1">
                        <CalendarIcon className="w-4 h-4 text-gray-500 mr-1" />
                        {formatJaDateSpan(confirmedDate?.start, confirmedDate?.end)}
                    </div>
                ) : (
                    <div className="space-y-1">
                        {event.proposed_dates?.slice(0, 2).map((date) => (
                            <div key={date.id} className="flex items-center text-sm text-gray-500">
                                <ChevronRightIcon className="w-4 h-4 text-yellow-500 mr-1" />
                                {formatJaDateSpan(date.start, date.end)}
                            </div>
                        ))}
                        {event.proposed_dates && event.proposed_dates.length > 2 && (
                            <p className="text-sm text-gray-400">... 他 {event.proposed_dates.length - 2} 件</p>
                        )}
                    </div>
                )}

            </div>
        </Card>
    );
};

export default EventCard;