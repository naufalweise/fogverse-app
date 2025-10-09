"use server";

import { BACKEND_URL } from "@/settings";
import { redirect } from 'next/navigation'
import { IGenericResponse } from "@/types/response";


export default async function deployContainer(data: Job): Promise<IGenericResponse> {
    let response: Response; // Declare response variable outside the try block

    try {
        response = await fetch(`${BACKEND_URL}/run/jobs`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });

    } catch (error) {
        console.error("Fetch failed:", error);
        if (error instanceof Error) {
            return { "status": "error", "message": error.message };
        }
        return { "status": "error", "message": "Error: Unable to connect to the server" };
    }

    // Handle API errors (like validation errors from the backend)
    if (!response.ok) {
        return await response.json();
    }

    // If we reach this point, the request was successful.
    // Now, we can safely redirect.
    redirect("/run/jobs");
}