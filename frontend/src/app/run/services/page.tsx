import { Container, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from "@mui/material";
import getServices from "./actions";



export default async function Page() {
    const services = await getServices();
    return <Container>
        <Typography variant="h4" component="h1">Services</Typography>
        <TableContainer component={Paper}>

        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Image</TableCell>
                    <TableCell>Input Topic</TableCell>
                    <TableCell>Output Topic</TableCell>
                    <TableCell>Replica</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {services.data?.map(s => (
                    <TableRow key={s.name}>
                        <TableCell>{s.name}</TableCell>
                        <TableCell>{s.readyReplicas > 0 ? "Ready" : "Not Ready"}</TableCell>
                        <TableCell>{s.imageURL}</TableCell>
                        <TableCell>{s.inputTopic}</TableCell>
                        <TableCell>{s.outputTopic}</TableCell>
                        <TableCell>{s.readyReplicas} / {s.replicas}</TableCell>
                    </TableRow>
                    ))}
            </TableBody>
        </Table>

        </TableContainer>
    </Container>
}