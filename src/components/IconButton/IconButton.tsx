'use client'
import React from 'react';
import { cva, cx, type VariantProps } from "class-variance-authority";
import { BaseButtonProps, BaseButton, baseButton } from '@/components/BaseButton/BaseButton';

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
        }
    },
    defaultVariants: {
        iconSize: "md",
        strokeWidth: 1,
    },
});

interface IconButtonProps extends BaseButtonProps, VariantProps<typeof icon> {};

const iconButton = ({ shape, intent, iconSize, strokeWidth }: IconButtonProps = {}) => cx(
    baseButton({ shape, intent }),
    icon({ iconSize, strokeWidth }),
);

const IconButton: React.FC<IconButtonProps> = ({ className, children, iconSize, strokeWidth, shape, intent, ...props }) => {
    return (
        <BaseButton
            className={iconButton({ shape, intent, iconSize, strokeWidth, className })}
            {...props}
        >
            {children}
        </BaseButton>
    )
}

export default IconButton;