import axios from "@/lib/axios/private";
import { AxiosError } from "axios";

export async function GET(request: Request) {
    return await axios.get('/auth/logout')
        .then(result => {
            const response = new Response(JSON.stringify(result.data), { status: result.status });

            const setCookieHeaders = result.headers['set-cookie']?.join(', ');
            if (setCookieHeaders) {
                response.headers.set('Set-Cookie', setCookieHeaders);
            }

            return response;
        })
        .catch((err: AxiosError) => {
            console.error(err);
            const status = err?.response?.status ? err.response.status : 500;
            return new Response(JSON.stringify(err?.response?.data), { status: status });
        });
}