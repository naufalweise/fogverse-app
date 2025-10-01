
import { BACKEND_URL } from "@/settings";
import { DataResponse } from "@/types/response";

export default async function getServices(): Promise<DataResponse<Service[]>> {
    let response: Response; // Declare response variable outside the try block

    try {
        response = await fetch(`${BACKEND_URL}/run/services`);

    } catch (error) {
        console.error("Fetch failed:", error);
        if (error instanceof Error) {
            return { "status": "error", "message": error.message };
        }
        return { "status": "error", "message": "Error: Unable to connect to the server" };
    }

    return await response.json();
}