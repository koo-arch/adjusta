import React from 'react';
import { PencilSquareIcon } from '@heroicons/react/20/solid';
import IconButton from '@/components/IconButton';

interface EditButtonProps {
    to: string;
}

const EditButton: React.FC<EditButtonProps> = ({ to }) => {
    return (
        <IconButton
            iconColor={'primary'}
            to={to}
        >
            <PencilSquareIcon />
        </IconButton>
    )
}

export default EditButton;