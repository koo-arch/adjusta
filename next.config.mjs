/** @type {import('next').NextConfig} */
const nextConfig = {
    async rewrites() {
        const backend = process.env.NEXT_PUBLIC_API_BASE_URL // .env.local や Vercel の環境変数で定義

        return [
            // Next.js 側で処理したい API ルート
            {
                source: '/api/auth/cookie',
                destination: '/api/auth/cookie',
            },
            // その他の /api/* はバックエンドへ
            {
                source: '/api/:path*',
                destination: `${backend}/api/:path*`,
            },
            // OAuth コールバック等の /auth/* もバックエンドへ
            {
                source: '/auth/:path*',
                destination: `${backend}/auth/:path*`,
            },
        ]
    },
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'tailwindui.com',
                port: '',
            },
            {
                protocol: 'https',
                hostname: 'images.unsplash.com',
                port: '',
            },
            {
                protocol: 'https',
                hostname: 'lh3.googleusercontent.com',
                port: '',
            }
        ]
    }
};

export default nextConfig;
