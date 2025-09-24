export default interface IDeployData {
    name: string;
    imageURL: string;
    type: string;
    replicas: number;
    inputTopic: string;
    outputTopic: string;
}
