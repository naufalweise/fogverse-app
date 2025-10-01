export interface IGenericResponse {
    message: string;
    status: string;
}

export interface DataResponse<Type> extends IGenericResponse {
    data?: Type
}