'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { cva, type VariantProps } from "class-variance-authority";

export const baseButton = cva("text-center", {
    variants: {
        variant: {
            outline: "border bg-transparent",
            solid: "",
        },
        shape: {
            rounded: "rounded",
            md: "rounded-md",
            lg: "rounded-lg",
            full: "rounded-full",
        },
        intent: { primary: "", secondary: "", danger: "", warning: "", success: "", clear: "" },
    },
    compoundVariants: [
        {
            variant: "outline",
            intent: "primary",
            className: "border-indigo-500 text-indigo-500"
        },
        {
            variant: "outline",
            intent: "secondary",
            className: "border-pink-500 text-pink-500"
        },
        {
            variant: "outline",
            intent: "danger",
            className: "border-red-500 text-red-500"
        },
        {
            variant: "outline",
            intent: "warning",
            className: "border-yellow-500 text-yellow-500"
        },
        {
            variant: "outline",
            intent: "success",
            className: "border-green-500 text-green-500"
        },
        {
            variant: "outline",
            intent: "clear",
            className: "border-inherit text-inherit hover:border-gray-500 hover:shadow-md active:border-gray-600 active:shadow-sm active:bg-inherit"
        },
        {
            variant: "solid",
            intent: "primary",
            className: "bg-indigo-500 hover:bg-indigo-600 active:bg-indigo-700 text-white"
        },
        {
            variant: "solid",
            intent: "secondary",
            className: "bg-pink-500 hover:bg-pink-600 active:bg-pink-700 text-white"
        },
        {
            variant: "solid",
            intent: "danger",
            className: "bg-red-500 hover:bg-red-600 active:bg-red-700 text-white"
        },
        {
            variant: "solid",
            intent: "warning",
            className: "bg-yellow-500 hover:bg-yellow-600 active:bg-yellow-700 text-white"
        },
        {
            variant: "solid",
            intent: "success",
            className: "bg-green-500 hover:bg-green-600 active:bg-green-700 text-white"
        },
        {
            variant: "solid",
            intent: "clear",
            className: "bg-transparent hover:opacity-60 active:opacity-20 text-inherit"
        }
    ],
    defaultVariants: {
        variant: "solid",
        shape: "rounded",
        intent: "primary",
    },
});

export interface BaseButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof baseButton> {
    to?: string;
}

export const BaseButton: React.FC<BaseButtonProps> = ({ className, onClick, to, children, ...props }) => {
    const router = useRouter();

    const handleClick = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        if (onClick) {
            onClick(event);
        }

        if (!event.defaultPrevented && to) {
            if (to.startsWith("http") || to.startsWith("https")) {
                window.location.href = to;
            } else {
                router.push(to);
            }
        }
    }

    return (
        <button
            className={className}
            onClick={handleClick}
            {...props}
        >
            {children}
        </button>
    )
}