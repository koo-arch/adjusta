'use client'
import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useLogout } from '@/hooks/auth/useLogout';
import { usePathname } from 'next/navigation';
import { Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/20/solid';

const navigation = [
    { name: 'Home', href: '/' },
    { name: 'About', href: '/about' },
    { name: 'Contact', href: '/contact' },
]

const classNames = (...classes: string[]) => {
    return classes.filter(Boolean).join(' ')
}

const Header: React.FC = () => {
    const pathname = usePathname();
    const handleLogout = useLogout();
    return (
        <Disclosure as="nav" className="lg:pb-12">
            {({ open }) => (
                <>
                    <div className="mx-auto max-w-screen-2xl px-4 md:px-8">
                        <div className='relative flex items-center justify-between py-4 md:py-8'>
                            <div className="absolute inset-y-0 left-0 flex items-center md:hidden">

                                <DisclosureButton className="relative inline-flex items-center justify-center rounded-md p-2">
                                    <span className="absolute -inset-0.5"></span>
                                    <span className="sr-only">Open main menu</span>
                                    {open ? (
                                        <XMarkIcon className="h-6 w-6" aria-hidden="true" />
                                    ) : (
                                        <Bars3Icon className="h-6 w-6" aria-hidden="true" />
                                    )}
                                </DisclosureButton>
                            </div>
                            <div className="flex flex-1 items-center justify-center md:items-stretch md:justify-start">
                                <div className="flex flex-shrink-0 items-center">
                                    <Image
                                        className="h-6 w-auto"
                                        height={24}
                                        width={24}
                                        src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=500"
                                        alt="Your Company"
                                    />
                                </div>
                                <div className="hidden md:ml-6 md:block">
                                    <div className="flex gap-12">
                                        {navigation.map((item) => (
                                            <Link key={item.name} href={item.href} className={classNames(
                                                item.href === pathname ? 'text-indigo-500' : 'hover:text-indigo-500 active:text-indigo-700',
                                                "text-lg font-semibold transiton duration-100"
                                            )}
                                            >
                                                    {item.name}
                                            </Link>
                                        ))}
                                    </div>
                                </div>
                            </div>
                            <div className="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
                                <Menu as="div" className="relative ml-3">
                                    <div>
                                        <MenuButton className="relative flex rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                                            <span className="absolute -inset-1.5" />
                                            <span className="sr-only">Open user menu</span>
                                            <Image
                                                className="h-8 w-8 rounded-full"
                                                src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
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
                                                href="#"
                                                className={classNames(focus ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700')}
                                                >
                                                    Your Profile
                                                </a>
                                            )}
                                        </MenuItem>
                                        <MenuItem>
                                            {({ focus }) => (
                                                <a
                                                href="#"
                                                className={classNames(focus ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700')}
                                                >
                                                    Settings
                                                </a>
                                            )}
                                        </MenuItem>
                                        <MenuItem>
                                            {({ focus }) => (
                                                <a
                                                className={classNames(focus ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700')}
                                                onClick={handleLogout}
                                                >
                                                    Sign out
                                                </a>
                                            )}
                                        </MenuItem>
                                    </MenuItems>
                                </Menu>
                            </div>
                        </div>
                    </div>

                    <DisclosurePanel 
                        transition
                        className="sm:hidden origin-top transition duration-200 ease-out data-[closed]:-translate-x-6 data-[closed]:opacity-0"
                    >
                        <div className="space-y-1 px-4 pb-3 pt-2">
                            {navigation.map((item) => (
                                <DisclosureButton
                                    key={item.name}
                                    as={Link}
                                    href={item.href}
                                    className={classNames(
                                        item.href === pathname ? 'text-indigo-500' : 'hover:text-indigo-500 active:text-indigo-700',
                                        "block text-lg font-semibold"
                                    )}
                                >
                                    {item.name}
                                </DisclosureButton>
                            ))}
                        </div>
                    </DisclosurePanel>
                </>
            )}
        </Disclosure>
    )
}

export default Header;