import type { PageLoad } from "./$types.js";
import { PUBLIC_BACKEND_BASE_URL } from "$env/static/public";


export const load: PageLoad = async ({}) => {
    const res = await fetch(`${PUBLIC_BACKEND_BASE_URL}/v1/stats`);
	const stats = await res.json();
    return { 
        stats: stats,
        title: "Summary"
    };
};