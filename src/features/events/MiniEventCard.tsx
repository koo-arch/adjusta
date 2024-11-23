import React from 'react';
import Card from '@/components/Card';
import { formatJaDateSpan } from '@/lib/date/format';
import { CalendarIcon, ExclamationCircleIcon } from '@heroicons/react/20/solid';

interface MiniEventCardProps {
    title : string;
    start: Date;
    end: Date;
    needs_attention?: boolean;
    onClick: () => void;
}

const MiniEventCard: React.FC<MiniEventCardProps> = ({ title, start, end, needs_attention, onClick }) => {
    return (
        <Card
            isButton
            variant="outlined"
            background="white"
            onClick={onClick}
        >
            <div className="flex justify-between items-center">
                <h3 className="text-indigo-600 font-semibold truncate">{title}</h3>
               {needs_attention && 
                   <ExclamationCircleIcon className="w-5 h-5 text-red-500" title="要調整" />
                }
            </div>
            <div className="flex items-center mt-2">
                <CalendarIcon className="w-4 h-4 text-gray-500 mr-2" />
                <p className="text-xs text-gray-700">{formatJaDateSpan(start, end)}</p>
            </div>
        </Card>
    )
}

export default MiniEventCard;