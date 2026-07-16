import React from 'react';

const MarketingFooter = async () => {
    'use cache';

    return (
        <footer className="border-t border-gray-200 bg-white">
            <div className="mx-auto max-w-screen-2xl px-4 py-6 text-center text-sm text-gray-500 md:px-8">
                &copy; {new Date().getFullYear()} Adjusta
            </div>
        </footer>
    );
};

export default MarketingFooter;
