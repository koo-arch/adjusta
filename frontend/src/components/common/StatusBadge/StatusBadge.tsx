import React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

const dotVariants = cva('rounded-full', {
    variants: {
        size: {
            sm: 'size-1',
            md: 'size-2',
            lg: 'size-3',
        },
        color: {
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
        size: 'md',
        color: 'gray',
    },
});

const labelVariants = cva('', {
    variants: {
        size: {
            sm: 'text-sm',
            md: 'text-base',
            lg: 'text-lg',
            xl: 'text-xl',
        },
        color: {
            gray: 'text-gray-500',
            red: 'text-red-500',
            green: 'text-green-500',
            blue: 'text-blue-500',
            yellow: 'text-yellow-600',
            indigo: 'text-indigo-500',
            purple: 'text-purple-500',
            pink: 'text-pink-500',
        },
    },
    defaultVariants: {
        size: 'md',
        color: 'gray',
    },
});

interface StatusBadgeProps {
    label: string;
    color?: VariantProps<typeof dotVariants>['color'];
    dotSize?: VariantProps<typeof dotVariants>['size'];
    textSize?: VariantProps<typeof labelVariants>['size'];
    className?: string;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({
    label,
    color,
    dotSize,
    textSize,
    className,
}) => (
    <Badge
        variant="outline"
        className={cn('gap-2 border-0 bg-transparent p-0 font-normal shadow-none', className)}
    >
        <span aria-hidden="true" className={dotVariants({ color, size: dotSize })} />
        <span className={labelVariants({ color, size: textSize })}>{label}</span>
    </Badge>
);

export default StatusBadge;
