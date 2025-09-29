"use server";

import { BACKEND_URL } from "@/settings";
import IDeployData from "../types/DeployData";
import { redirect } from 'next/navigation'
import { IGenericResponse } from "@/types/response";


export default async function deployContainer(data: IDeployData): Promise<IGenericResponse> {
    console.log("ceck")
    try {
        const response = await fetch(`${BACKEND_URL}/deploy`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
        if (response.ok) {
            redirect("/deploy")
        }
        return await response.json();

    } catch (error) {
        console.error(error)
        if (error instanceof Error) {
            return {"status": "error", "message": error.message}
        }
        return {"status": "error", "message": "Error: Unable to process your request" }
    }
    
}