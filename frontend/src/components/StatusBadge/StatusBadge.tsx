import React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

const badgeStyle = cva('rounded-full', {
    variants: {
        circleSize: {
            sm: 'p-0.5',
            md: 'p-1',
            lg: 'p-1.5',
        },
        circleColor: {
            gray: 'bg-gray-500',
            red: 'bg-red-500',
            green: 'bg-green-500',
            blue: 'bg-blue-500',
            yellow: 'bg-yellow-500',
            indigo: 'bg-indigo-500',
            purple: 'bg-purple-500',
            pink: 'bg-pink-500',
        },
    },
    defaultVariants: {
        circleSize: 'md',
        circleColor: 'gray',
    },
})

const textStyle = cva("", {
    variants: {
        textSize: {
            sm: "text-sm",
            md: "text-base",
            lg: "text-lg",
            xl: "text-xl",
        },
        textColor: {
            gray: "text-gray-500",
            red: "text-red-500",
            green: "text-green-500",
            blue: "text-blue-500",
            yellow: "text-yellow-500",
            indigo: "text-indigo-500",
            purple: "text-purple-500",
            pink: "text-pink-500",
        }
    },
    defaultVariants: {
        textSize: "md",
        textColor: "gray",
    },
});

interface StatusBadgeProps extends VariantProps<typeof badgeStyle>, VariantProps<typeof textStyle> {
    label: string;
    className?: string;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ label, circleColor, circleSize, textColor, textSize, className }) => {
    return (
        <div className="flex items-center space-x-2">
            <span className={`${badgeStyle({ circleColor, circleSize })} ${className}`}>
            </span>
            <span className={`${textStyle({ textColor, textSize })}`}>{label}</span>
        </div>
    );
};

export default StatusBadge;