'use client'
import React from 'react';
import IconButton from '@/components/IconButton';
import { PlusIcon } from '@heroicons/react/20/solid';

const DraftRegisterButton = () => {
    return (
        <IconButton
            iconColor="primary"
            iconSize="lg"
            to="/schedule/draft/register"
        >
            <PlusIcon />
        </IconButton>
    )
}

export default DraftRegisterButton;
