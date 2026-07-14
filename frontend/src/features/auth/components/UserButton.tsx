'use client'
import React from 'react';
import Image from 'next/image';
import Link from 'next/link';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useLogout } from '@/features/auth/hooks/useLogout';
import type { AuthUser } from '@/features/auth/types';
import { LogOut, User } from 'lucide-react';

interface UserButtonProps {
    user: AuthUser;
}

const UserButton: React.FC<UserButtonProps> = ({ user }) => {
    const { logout } = useLogout();

    return (
        <DropdownMenu>
            <DropdownMenuTrigger className="rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
                <span className="sr-only">ユーザーメニューを開く</span>
                {user.picture ? (
                    <Image
                        className="h-8 w-8 rounded-full"
                        src={user.picture}
                        alt={`${user.name}のプロフィール画像`}
                        width={32}
                        height={32}
                    />
                ) : (
                    <span className="grid size-8 place-items-center rounded-full bg-muted">
                        <User className="size-4" aria-hidden="true" />
                    </span>
                )}
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
                <DropdownMenuItem asChild>
                    <Link href="/account">
                        <User />
                        アカウント
                    </Link>
                </DropdownMenuItem>
                <DropdownMenuItem
                    onSelect={() => {
                        void logout();
                    }}
                >
                    <LogOut />
                    ログアウト
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}

export default UserButton;
