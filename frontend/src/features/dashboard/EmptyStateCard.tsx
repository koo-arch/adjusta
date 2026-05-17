import React from 'react';
import Card from '@/components/Card';

interface EmptyStateCardProps {
    children: React.ReactNode;
}

const EmptyStateCard: React.FC<EmptyStateCardProps> = ({ children }) => {
    return (
        <Card
            variant="outlined"
            background="inherit"
            className="p-6 flex justify-center items-center text-xl font-bold text-gray-700 dark:text-gray-300"
        >
            {children}
        </Card>
    )
}

export default EmptyStateCard;