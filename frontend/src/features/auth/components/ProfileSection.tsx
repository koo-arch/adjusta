'use client'
import React, { useState } from 'react';
import Image from 'next/image';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import type { AuthUser } from '@/features/auth/types';

interface ProfileSectionProps {
    user: AuthUser;
}

const ProfileSection: React.FC<ProfileSectionProps> = ({ user }) => {
    const [imageFailed, setImageFailed] = useState(false);
    const showFallback = !user.picture || imageFailed;

    return (
        <Card>
            <CardHeader>
                <CardTitle>プロフィール</CardTitle>
                <CardDescription>
                    Google アカウントの情報のため、このアプリからは変更できません。
                </CardDescription>
            </CardHeader>
            <CardContent className="flex items-center gap-4">
                {showFallback ? (
                    <div
                        aria-hidden
                        className="flex h-16 w-16 shrink-0 items-center justify-center rounded-full bg-primary text-2xl font-bold text-primary-foreground"
                    >
                        {user.name.charAt(0)}
                    </div>
                ) : (
                    <Image
                        className="h-16 w-16 shrink-0 rounded-full"
                        src={user.picture}
                        width={64}
                        height={64}
                        alt=""
                        onError={() => setImageFailed(true)}
                    />
                )}
                <div className="min-w-0">
                    <p className="truncate text-base font-medium">{user.name}</p>
                    <p className="truncate text-sm text-muted-foreground">{user.email}</p>
                </div>
            </CardContent>
        </Card>
    );
};

export default ProfileSection;
