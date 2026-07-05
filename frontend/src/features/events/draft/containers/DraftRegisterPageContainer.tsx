import React from 'react';
import EventDraft from '@/features/events/draft/components/EventDraft';

const DraftRegisterPageContainer = () => {
    return (
        <div className="mx-auto p-4 max-w-screen-lg">
            <h1 className="text-3xl font-bold text-center mb-8">イベント登録</h1>
            <EventDraft />
        </div>
    );
};

export default DraftRegisterPageContainer;
