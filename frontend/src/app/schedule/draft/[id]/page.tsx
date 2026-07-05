import React from 'react';
import EventDetailPageContainer from '@/features/events/detail/containers/EventDetailPageContainer';

interface DraftDetailPageProps {
    params: {
        id: string;
    };
}

const DraftDetailPage = ({ params }: DraftDetailPageProps) => {
    return <EventDetailPageContainer eventID={params.id} />;
};

export default DraftDetailPage;
