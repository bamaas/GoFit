import { toast } from "svelte-sonner";
import { getAuthToken } from "./auth";
import { goto } from "$app/navigation";
import { deleteCookie } from "./cookie";

function catchError(error: Error){
    if (error.message === "Unauthorized"){
        deleteCookie("token");
        toast.error("You are not authorized to view this page.", {
            description: "Please login.", 
            cancel: { label: "X" }
        });
        goto("/login")
    } else {
        console.log(error)
        toast.error("Something went wrong.", {
            description: "Oops!", 
            cancel: { label: "X" }
        });  
    }
}

export async function request(url: string, options?: any): Promise<any> {
    const authToken = getAuthToken();

    if (options === undefined) {
        options = {
            headers: {}
        };
    }
    options["headers"] = {}
    // options.headers["Content-Type"] = "application/json";
    options.headers["Authorization"] = `Bearer ${authToken}`;

    return fetch(url, options)
    .then( response => {
        if (!response.ok){
            catchError(new Error(response.statusText))
        }
        return response
    })
    .then((response) => response.json())
    .catch(catchError)
}