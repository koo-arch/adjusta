import React from 'react';
import EventList from '@/features/events/EventList';

const DraftPage = () => {
    return (
        <div>
            <h1>Draft</h1>
            <div className="max-w-screen-md mx-auto px-4">
                <EventList />
            </div>
        </div>
    )
}

export default DraftPage;