'use client'
import React from 'react';
import { cva, cx, type VariantProps } from "class-variance-authority";
import { BaseButtonProps, BaseButton, baseButton } from '../BaseButton/BaseButton';

const buttonSize = cva("inline-flex items-center justify-center", {
    variants: {
        size: {
            sm: "px-2 py-1 text-sm",
            md: "px-4 py-2 text-base",
            lg: "px-6 py-3 text-lg",
        },
    },
    defaultVariants: {
        size: "md",
    },
});

const icon = cva("", {
    variants: {
        iconSize: {
            sm: "h-4 w-4",
            md: "h-5 w-5",
            lg: "h-6 w-6",
        },
    },
    defaultVariants: {
        iconSize: "md",
    },
});


interface ButtonProps extends BaseButtonProps, VariantProps<typeof icon>, VariantProps<typeof buttonSize> {
    startIcon?: React.ReactNode;
    endIcon?: React.ReactNode;
}

const button = ({ shape, intent, variant, size }: ButtonProps = {}) => cx(
    baseButton({ shape, intent, variant }),
    buttonSize({ size }),
);

const Button: React.FC<ButtonProps> = ({ children, startIcon, endIcon, iconSize, size, className, ...props }) => {
    const { shape, intent, variant } = props;

    return (
        <BaseButton
            className={button({ shape, intent, size, variant, className })}
            {...props}
        >
            {startIcon && <span className={`${icon({ iconSize })} mr-2`}>{startIcon}</span>}
            {children}
            {endIcon && <span className={`${icon({ iconSize })} ml-2`}>{endIcon}</span>}
        </BaseButton>
    )
}

export default Button;