import React from 'react';
import Image from 'next/image';
import type { AuthUser } from '@/features/auth/types';

interface UserInfoProps {
    user: AuthUser;
}

const UserInfo: React.FC<UserInfoProps> = ({ user }) => {
    return (
        <div>
            <div className="ma-auto grid grid-cols-3 gap-4 mb-4">
                <section className="col-span-1 flex justify-center items-center">
                    <Image
                        className="rounded-full"
                        src={user.picture}
                        width={100}
                        height={100}
                        alt={user.name}
                    />
                </section>
                <section className="col-span-2">
                    <p className="font-bold text-2xl mb-2">{user.name}</p>
                    <p className="text-gray-500">{user.email}</p>
                </section>
            </div>
            <div className="border-t border-gray-300 mb-4"></div>
        </div>
    )
}

export default UserInfo;
