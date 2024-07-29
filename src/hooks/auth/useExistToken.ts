import axios from '@/lib/axios/middle';
import useSWR from 'swr';

const fetcher = async (url: string) => await axios.get(url).then(res => res.data);

export const useExistToken = () => {
    const { data, error } = useSWR(
        'api/auth/cookie',
        fetcher
    );

    const isExistToken: boolean = data?.exist

    return {
        isExistToken: isExistToken,
        error
    }
};