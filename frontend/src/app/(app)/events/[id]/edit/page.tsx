import React from 'react';
import EventEditPageContainer from '@/features/events/edit/containers/EventEditPageContainer';

interface EventEditPageProps {
    params: {
        id: string;
    };
}

const EventEditPage = ({ params }: EventEditPageProps) => {
    return <EventEditPageContainer eventID={params.id} />;
}

export default EventEditPage;
