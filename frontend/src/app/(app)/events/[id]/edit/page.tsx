import React from 'react';
import EventEditPageContainer from '@/features/events/edit/containers/EventEditPageContainer';

interface EventEditPageProps {
    // Next.js 15+ では params は Promise
    params: Promise<{
        id: string;
    }>;
}

const EventEditPage = async ({ params }: EventEditPageProps) => {
    const { id } = await params;
    return <EventEditPageContainer eventID={id} />;
}

export default EventEditPage;
