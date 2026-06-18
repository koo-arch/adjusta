import React from 'react';
import EventDetail from '@/features/events/detail/components/EventDetail';

interface EventDetailPageContainerProps {
    eventID: string;
}

const EventDetailPageContainer: React.FC<EventDetailPageContainerProps> = ({ eventID }) => {
    return (
        <div className="px-4 max-w-screen-lg mx-auto">
            <EventDetail eventID={eventID} />
        </div>
    );
};

export default EventDetailPageContainer;
