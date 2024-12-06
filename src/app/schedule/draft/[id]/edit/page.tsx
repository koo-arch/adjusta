import React from 'react';
import EventEdit from '@/features/events/edit/EventEdit';

const EditPage = () => {
    return (
        <div className="mx-auto p-4 max-w-screen-lg">
            <h1 className="text-3xl font-bold text-center mb-8">イベント編集</h1>
            <EventEdit />
        </div>
    )
}

export default EditPage;