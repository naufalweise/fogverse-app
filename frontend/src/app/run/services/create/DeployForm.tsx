"use client";

import { TextField, Button, Box, Alert, AlertTitle, FormControlLabel } from '@mui/material';
import { SubmitHandler, useForm } from 'react-hook-form'
import deployContainer from './actions';
import Checkbox from '@mui/material/Checkbox';
import InputAdornment from '@mui/material/InputAdornment';

type Props = {}

// TODO: Implement validation on server on topic must not be the same
function DeployForm() {
    const { register, handleSubmit, setError, watch, formState: {isLoading, errors}} = useForm<Service>({
        defaultValues: {
            replicas: 1,
            kedaConfig: {
                minReplicaCount: 0,
                maxReplicaCount: 10,
                cooldownPeriod: 5,
                pollingInterval: 10,
                lagThreshold: 1,

            }
        }
    });
    const onSubmit: SubmitHandler<Service> = async(data) => {
        console.log(data)
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
            }}
            noValidate
            autoComplete="off"
            onSubmit={handleSubmit(onSubmit)}
            >
            <TextField label="Name" variant='outlined' {...register("name")} margin='normal'/>
            <TextField label="Container Image URL" variant='outlined' {...register("imageURL")} margin='normal'/>
            <TextField label="Input Topic" {...register("inputTopic")} margin='normal' />
            <TextField label="Output Topic" {...register("outputTopic")} margin='normal' />
            <TextField type='number' label="Replicas" {...register("replicas", {valueAsNumber: true})} margin='normal'/>
            <FormControlLabel  control={<Checkbox slotProps={{input: {
                ...register("useKeda")
            }}} />} label="Use Keda" />
            {watch("useKeda") && (
                <>
                 <TextField 
                    label = "Min Replica Count"
                    {...register("kedaConfig.minReplicaCount", {valueAsNumber: true})}
                    helperText="Minimum number of replicas KEDA will scale the resource down to"
                    margin='normal'
                    type='number'
                    />
                 <TextField
                    label = "Max Replica Count"
                    {...register("kedaConfig.maxReplicaCount", {valueAsNumber: true})}
                    helperText="This setting is passed to the HPA definition that KEDA will create for a given resource and holds the maximum number of replicas of the target resource."
                    margin='normal'
                    type='number'
                    />
                <TextField 
                    label = "Lag Threshold"
                    {...register("kedaConfig.lagThreshold", {valueAsNumber: true})}
                    helperText="Target value for the total lag (sum of all partition lags) to trigger scaling actions."
                    margin='normal'
                    type='number'
                    />
                 <TextField 
                    label = "Cooldown Period"
                    {...register("kedaConfig.cooldownPeriod", {valueAsNumber: true})}
                    helperText="The period to wait after the last trigger reported active before scaling the resource back to 0, in seconds."
                    slotProps={{
                        input: {
                            endAdornment: <InputAdornment position="end">seconds</InputAdornment>,
                        },
                    }}
                    margin='normal'
                    type='number'
                    />

                 <TextField 
                    label = "Polling Interval"
                    {...register("kedaConfig.pollingInterval", {valueAsNumber: true})}
                    helperText="Interval to check each trigger on, in seconds"
                    slotProps={{
                        input: {
                            endAdornment: <InputAdornment position="end">seconds</InputAdornment>,
                        },
                    }}
                    type='number'
                    margin='normal'
                    />
                 
                </>
            )}
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