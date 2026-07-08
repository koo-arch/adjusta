import React from 'react';
import EventEditPageContainer from '@/features/events/edit/containers/EventEditPageContainer';

interface EditPageProps {
    params: {
        id: string;
    };
}

const EditPage = ({ params }: EditPageProps) => {
    return <EventEditPageContainer eventID={params.id} />;
}

export default EditPage;
