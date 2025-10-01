import { Container, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Button } from "@mui/material";
import getServices from "./actions";
import Link from "next/link";



export default async function Page() {
    const services = await getServices();
    return <Container>
        <Typography variant="h4" component="h1" mb="1rem">Services</Typography>
        <Link href="/deploy">
            <Button sx={{mb: "1rem"}}>Create Services</Button>
        </Link>
        
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