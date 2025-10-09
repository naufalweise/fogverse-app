type Job = {
    name: string;
    imageURL: string;
    inputTopic: string;
    outputTopic: string;
    completion: number;
    parallelism: number;
    status: {
        succeeded: number;
    }
}