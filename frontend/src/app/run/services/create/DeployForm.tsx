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
                maxReplicaCount: 100,
                cooldownPeriod: 300,
                pollingInterval: 30,
                lagThreshold: 10,

            }
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
            {watch("useKeda") && (
                <>
                 <TextField 
                    label = "Min Replica Count"
                    {...register("kedaConfig.minReplicaCount")}
                    helperText="Minimum number of replicas KEDA will scale the resource down to"
                    />
                 <TextField
                    label = "Max Replica Count"
                    {...register("kedaConfig.maxReplicaCount")}
                    helperText="This setting is passed to the HPA definition that KEDA will create for a given resource and holds the maximum number of replicas of the target resource."
                    />
                 <TextField 
                    label = "Cooldown Period"
                    {...register("kedaConfig.cooldownPeriod")}
                    helperText="The period to wait after the last trigger reported active before scaling the resource back to 0, in seconds."
                    slotProps={{
                        input: {
                            endAdornment: <InputAdornment position="end">seconds</InputAdornment>,
                        },
                    }}
                    />

                 <TextField 
                    label = "Polling Interval"
                    {...register("kedaConfig.pollingInterval")}
                    helperText="Interval to check each trigger on, in seconds"
                    slotProps={{
                        input: {
                            endAdornment: <InputAdornment position="end">seconds</InputAdornment>,
                        },
                    }}
                    />
                 <TextField 
                    label = "Lag Threshold"
                    {...register("kedaConfig.lagThreshold")}
                    helperText="Target value for the total lag (sum of all partition lags) to trigger scaling actions."
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