'use client'
import React from 'react';
import { useTheme } from 'next-themes';
import IconButton from '@/components/IconButton';
import { MoonIcon, SunIcon } from '@heroicons/react/20/solid';
import Image from 'next/image';

interface ThemeButtonProps {
    className?: string;
}

const ThemeButton: React.FC<ThemeButtonProps> = ({ className }) => {
    const { resolvedTheme, setTheme } = useTheme();

    if (!resolvedTheme) return (
        <Image
            src="data:image/svg+xml;base64,PHN2ZyBzdHJva2U9IiNGRkZGRkYiIGZpbGw9IiNGRkZGRkYiIHN0cm9rZS13aWR0aD0iMCIgdmlld0JveD0iMCAwIDI0IDI0IiBoZWlnaHQ9IjIwMHB4IiB3aWR0aD0iMjAwcHgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHJlY3Qgd2lkdGg9IjIwIiBoZWlnaHQ9IjIwIiB4PSIyIiB5PSIyIiBmaWxsPSJub25lIiBzdHJva2Utd2lkdGg9IjIiIHJ4PSIyIj48L3JlY3Q+PC9zdmc+Cg=="
            width={36}
            height={36}
            sizes="36x36"
            alt="Loading Light/Dark Toggle"
            priority={false}
            title="Loading Light/Dark Toggle"
        />
    )

    return (
        <IconButton
            iconColor="clear"
            iconSize='lg'
            className={className}
            onClick={() => setTheme(resolvedTheme === 'dark' ? 'light' : 'dark')}
        >
            {resolvedTheme === 'light' ? <SunIcon /> : <MoonIcon />}
        </IconButton>
    )
}

export default ThemeButton;
