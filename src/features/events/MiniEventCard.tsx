import React from 'react';
import Card from '@/components/Card';
import { UpcomingEvent } from '@/hooks/event/type';
import { formatJaDateSpan } from '@/lib/date/format';
import { MapPinIcon } from '@heroicons/react/20/solid';

interface MiniEventCardProps {
    event: UpcomingEvent;
    onClick: () => void;
}

const MiniEventCard: React.FC<MiniEventCardProps> = ({ event, onClick }) => {
    return (
        <Card
            isButton
            variant="outlined"
            background="white"
            onClick={onClick}
        >
            <h3 className="text-indigo-600 font-semibold truncate">{event.title}</h3>
            <p className="text-sm text-gray-500">
                {formatJaDateSpan(event.start, event.end)}
            </p>
            <div className="flex items-center mt-2">
                <MapPinIcon className="w-4 h-4 text-gray-500 mr-2" />
                <p className="text-xs text-gray-700">{event.location || '未設定'}</p>
            </div>
        </Card>
    )
}

export default MiniEventCard;