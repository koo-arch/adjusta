import React from 'react';
import EventDetailPageContainer from '@/features/events/detail/containers/EventDetailPageContainer';

interface EventDetailPageProps {
    // Next.js 15+ では params は Promise
    params: Promise<{
        id: string;
    }>;
}

const EventDetailPage = async ({ params }: EventDetailPageProps) => {
    const { id } = await params;
    return <EventDetailPageContainer eventID={id} />;
};

export default EventDetailPage;
