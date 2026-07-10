import React from 'react';
import EventDetailPageContainer from '@/features/events/detail/containers/EventDetailPageContainer';

interface EventDetailPageProps {
    params: {
        id: string;
    };
}

const EventDetailPage = ({ params }: EventDetailPageProps) => {
    return <EventDetailPageContainer eventID={params.id} />;
};

export default EventDetailPage;
