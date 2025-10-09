import { Container, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Button, Alert, } from "@mui/material";
import getJobs from "./actions";
import Link from "next/link";



export default async function Page() {
    const response = await getJobs();
    return <Container>
        <Typography variant="h4" component="h1" mb="1rem">Jobs</Typography>
        <Link href="/run/jobs/create">
            <Button sx={{mb: "1rem"}}>Create Job</Button>
        </Link>
        
        {response.status == "error" && <Alert severity="error">{response.message}</Alert>}
        <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Completions</TableCell>
                    <TableCell>Image</TableCell>
                    <TableCell>Input Topic</TableCell>
                    <TableCell>Output Topic</TableCell>
                    <TableCell>Paralellism</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {response.data?.map(s => (
                    <TableRow key={s.name}>
                        <TableCell>{s.name}</TableCell>
                        <TableCell>{s.status.succeeded} / {s.completion} </TableCell>
                        <TableCell>{s.imageURL}</TableCell>
                        <TableCell>{s.inputTopic}</TableCell>
                        <TableCell>{s.outputTopic}</TableCell>
                        <TableCell>{s.parallelism}</TableCell>
                    </TableRow>
                    ))}
            </TableBody>
        </Table>

        </TableContainer>
    </Container>
}