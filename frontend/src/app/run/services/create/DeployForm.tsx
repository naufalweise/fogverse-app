"use client";

import { TextField, Button, Box, Alert, AlertTitle, FormControlLabel } from '@mui/material';
import { SubmitHandler, useForm } from 'react-hook-form'
import deployContainer from './actions';
import Checkbox from '@mui/material/Checkbox';

type Props = {}

// TODO: Implement validation on server on topic must not be the same
function DeployForm() {
    const { register, handleSubmit, setError,formState: {isLoading, errors}} = useForm<Service>({
        defaultValues: {
            replicas: 1
        }
    });
    const onSubmit: SubmitHandler<Service> = async(data) => {
        // console.log(data)
        const result = await deployContainer(data);
        if (result.status == "error") {
            setError("root", {"message": result.message})
        }
    }
    return (
        <Box
            component="form"
            sx={{
                display: 'flex',
                flexDirection: 'column',
                gap: "1em",
            }}
            noValidate
            autoComplete="off"
            onSubmit={handleSubmit(onSubmit)}
            >
            <TextField label="Name" variant='outlined' {...register("name")} />
            <TextField label="Container Image URL" variant='outlined' {...register("imageURL")} />
            <TextField label="Input Topic" {...register("inputTopic")} />
            <TextField label="Output Topic" {...register("outputTopic")} />
            <TextField type='number' label="Replicas" {...register("replicas")} />
            <FormControlLabel control={<Checkbox slotProps={{input: {
                ...register("useKeda")
            }}} />} label="Use Keda" />
            <Button variant="contained" type="submit" sx={{ mt: 2 }} disabled={isLoading}>
                Submit
            </Button>
            {errors.root && (
                <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {errors.root.message}
                </Alert>
            )}
        </Box>
  );
}

export default DeployForm;