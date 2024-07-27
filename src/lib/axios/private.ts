import axios, { AxiosInstance } from 'axios';
import { axiosConfig } from './public';
import { cookies } from 'next/headers';

const instance: AxiosInstance = axios.create(axiosConfig);

instance.interceptors.request.use((config) => {
    const accessToken = cookies().get('access_token');
    console.log(accessToken);

    if (accessToken) {
        config.headers.Authorization = `Bearer ${ accessToken.value }`;
    }

    config.headers['Cookie'] = cookies().toString();

    return config;
}, (error) => {
    return Promise.reject(error);
});

export default instance;