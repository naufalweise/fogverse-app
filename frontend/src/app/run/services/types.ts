type Service = {
    name: string;
    imageURL: string;
    inputTopic: string;
    outputTopic: string;
    replicas: number;
    readyReplicas: number;
    useKeda: boolean;
}

type KedaConfig = {

}