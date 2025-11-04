type Service = {
    name: string;
    imageURL: string;
    inputTopic: string;
    outputTopic: string;
    replicas: number;
    readyReplicas: number;
    useKeda: boolean;
    kedaConfig?: KedaConfig
}

type KedaConfig = {
    minReplicaCount: number;
    maxReplicaCount: number;
    cooldownPeriod: number;
    pollingInterval: number;
    lagThreshold: number;
}