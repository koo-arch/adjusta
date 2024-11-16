import React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

const cardStyle = cva('rounded-lg p-4', {
    variants: {
        variant: {
            shadow: 'shadow-lg',
            outlined: 'border border-gray-300',
        },
        background: {
            inherit: '',
            white: 'bg-white',
            gray: 'bg-gray-100',
        },
        isButton: {
            true: 'cursor-pointer',
            false: '',
        }
    },
    compoundVariants: [
        {
            variant: 'shadow',
            background: 'inherit',
            isButton: true,
            className: 'hover:opacity-60 hover:shadow-sm active:shadow-none',

        },
        {
            variant: 'shadow',
            background: 'white',
            isButton: true,
            className: 'shadow-lg hover:opacity-60',
        },
        {
            variant: 'shadow',
            background: 'gray',
            isButton: true,
            className: 'shadow-lg hover:opacity-60',
        },
        {
            variant: 'outlined',
            background: 'white',
            isButton: true,
            className: 'border border-gray-300 hover:border-indigo-500 hover:shadow-md active:shadow-none active:border-gray-300 active:bg-gray-50',
        },
        {
            variant: 'outlined',
            background: 'gray',
            isButton: true,
            className: 'border border-gray-300 hover:border-indigo-500 hover:shadow-md active:shadow-none active:border-gray-300 active:bg-gray-50',
        },
        {
            variant: 'outlined',
            background: 'inherit',
            isButton: true,
            className: 'hover:border-indigo-500 hover:shadow-md active:shadow-sm active:border-indigo-600 active:bg-gray-50',
        },
    ],
    defaultVariants: {
        variant: 'shadow',
        background: 'inherit',
    },

});

interface CardProps extends VariantProps<typeof cardStyle> {
    children: React.ReactNode;
    className?: string;
    footer?: React.ReactNode;
    actions?: React.ReactNode;
    onClick?: () => void;
}

const Card: React.FC<CardProps> = ({ variant, background, isButton, footer, actions, children, onClick, className }) => {
    return (
        <div className={`${cardStyle({ variant, background, isButton })} ${className}`} onClick={onClick}>
            {children}
            {actions && (
                <div className="flex justify-end space-x-2">
                    {actions}
                </div>
            )}
            {footer && (
                <div className="mt-4 pt-4 border-t border-gray-200">
                    {footer}
                </div>
            )}
        </div>
    );
};

export default Card;