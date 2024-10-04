'use client'
import React from 'react';
import Card from '@/components/Card';
import StatusBadge from '@/components/StatusBadge';
import { formatJaDate } from '@/lib/date/format';
import { EventDraftDetail } from '@/hooks/event/type';

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
    return (
        <Card variant="outlined" background="inherit" isButton={true} onClick={onClick}>
            <div>
                <div className="flex justify-between items-center mb-2">
                    <h2 className="text-lg font-bold">{event?.title}</h2>
                    <StatusBadge label={statusLabel(event.status)} color={statusColor(event.status)} />
                </div>
                {event.proposed_dates?.map((date) => (
                    <p key={date.id} className="text-sm text-gray-500 mb-1">
                        第{date.priority}希望：{formatJaDate(date.start)} 〜 {formatJaDate(date.end)}
                    </p>
                ))}
                <p className="text-sm mt-2">{event.description}</p>
            </div>
        </Card>
    );
};

export default EventCard;