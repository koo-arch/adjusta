import React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

const badgeStyle = cva('rounded-full', {
    variants: {
        size: {
            sm: 'px-2 py-1 text-xs',
            md: 'px-3 py-1 text-sm',
            lg: 'px-4 py-2 text-base',
        },
        color: {
            gray: 'bg-gray-500 text-white',
            red: 'bg-red-500 text-white',
            green: 'bg-green-500 text-white',
            blue: 'bg-blue-500 text-white',
            yellow: 'bg-yellow-500 text-white',
            indigo: 'bg-indigo-500 text-white',
            purple: 'bg-purple-500 text-white',
            pink: 'bg-pink-500 text-white',
        },
    },
    defaultVariants: {
        size: 'md',
        color: 'gray',
    },
    }
)

interface StatusBadgeProps extends VariantProps<typeof badgeStyle> {
    label: string;
    className?: string;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ label, color, size, className }) => {
    return (
        <span className={`${badgeStyle({ color, size })} ${className}`}>
            {label}
        </span>
    );
};

export default StatusBadge;