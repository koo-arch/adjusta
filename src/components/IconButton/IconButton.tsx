'use client'
import React from 'react';
import { cva, type VariantProps } from "class-variance-authority";
import { BaseButton } from '@/components/BaseButton/BaseButton';

const icon = cva("p-1", {
    variants: {
        iconSize: {
            sm: "h-6 w-6",
            md: "h-8 w-8",
            lg: "h-10 w-10",
        },
        strokeWidth: {
            1: "stroke-1",
            2: "stroke-2",
            3: "stroke-3",
        },
        iconColor: {
            primary: "text-indigo-500 hover:text-indigo-600 active:text-indigo-700 dark:text-white dark:hover:text-indigo-400 dark:active:text-indigo-500",
            success: "text-green-500 hover:text-green-600 active:text-green-700",
            warning: "text-yellow-500 hover:text-yellow-600 active:text-yellow-700",
            danger: "text-red-500 hover:text-red-600 active:text-red-700",
            clear: "text-inherit hover:text-gray-500 active:text-gray-600 dark:hover:text-gray-300 dark:active:text-gray-400",
        }
    },
    defaultVariants: {
        iconSize: "md",
        strokeWidth: 1,
        iconColor: "clear",
    },
});

interface IconButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof icon> {
    to?: string;
};


const IconButton: React.FC<IconButtonProps> = ({ className, children, iconSize, iconColor, strokeWidth, ...props }) => {
    return (
        <BaseButton
            className={icon({ iconSize, iconColor, strokeWidth, className })}
            {...props}
        >
            {children}
        </BaseButton>
    )
}

export default IconButton;