'use client'
import React from 'react';
import Image from 'next/image';
import { Menu, MenuButton, MenuItems, MenuItem } from '@headlessui/react';
import { useLogout } from '@/hooks/auth/useLogout';
import { useAuth } from '@/hooks/auth/useAuth';
import { UserIcon, ArrowRightStartOnRectangleIcon } from '@heroicons/react/20/solid';

interface UserButtonProps {
    classNames: (...classes: string[]) => string;
}

const UserButton: React.FC<UserButtonProps> = ({ classNames }) => {
    const handleLogout = useLogout();
    const { isAuthenticated, user, isLoading } = useAuth();

    if (isLoading) return null;
    if (!isAuthenticated || !user) return null;

    return (
        <div className="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
            <Menu as="div" className="relative ml-3">
                <div>
                    <MenuButton className="relative flex rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                        <span className="absolute -inset-1.5" />
                        <span className="sr-only">Open user menu</span>
                        <Image
                            className="h-8 w-8 rounded-full"
                            src={user?.picture}
                            alt=""
                            width={32}
                            height={32}
                        />
                    </MenuButton>
                </div>
                <MenuItems
                    transition
                    className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 transition focus:outline-none data-[closed]:scale-95 data-[closed]:transform data-[closed]:opacity-0 data-[enter]:duration-100 data-[leave]:duration-75 data-[enter]:ease-out data-[leave]:ease-in"
                >
                    <MenuItem>
                        {({ focus }) => (
                            <a
                                href="/account"
                                className={classNames(focus ? 'bg-gray-100' : '', 'flex items-center px-4 py-2 text-sm text-gray-700')}
                            >
                                <UserIcon className="h-5 w-5 mr-2 text-gray-500" />
                                アカウント
                            </a>
                        )}
                    </MenuItem>
                    <MenuItem>
                        {({ focus }) => (
                            <a
                                className={classNames(focus ? 'bg-gray-100' : '', 'flex items-center px-4 py-2 text-sm text-gray-700')}
                                onClick={handleLogout}
                            >
                                <ArrowRightStartOnRectangleIcon className="h-5 w-5 mr-2 text-gray-500" />
                                ログアウト
                            </a>
                        )}
                    </MenuItem>
                </MenuItems>
            </Menu>
        </div>
    )
}

export default UserButton;