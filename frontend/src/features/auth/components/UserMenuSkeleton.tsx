import React from 'react';

const UserMenuSkeleton: React.FC = () => {
    return (
        <div className="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
            <div className="ml-3 h-8 w-8 animate-pulse rounded-full bg-gray-200" />
        </div>
    );
};

export default UserMenuSkeleton;
