"use client";

import React, { useActionState, useState } from 'react';
import { TextField, Button, Box, Typography, FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, Alert, AlertTitle } from '@mui/material';
import { Controller, SubmitHandler, useForm } from 'react-hook-form'
import IDeployData from './types/DeployData';
import deployContainer from './actions/deployContainer';
import { IGenericResponse } from '@/types/response';


type Props = {}

// TODO: Implement validation on server on topic must not be the same
function DeployForm() {
    const { register, handleSubmit, control, setError,formState: {isLoading, errors}} = useForm<IDeployData>({
        defaultValues: {
            type: "service"
        }
    });
    const onSubmit: SubmitHandler<IDeployData> = async(data) => {
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
            <FormControl>
                <FormLabel>Type</FormLabel>
                <Controller control={control} name='type' render={({field}) => (
                    <RadioGroup {...field}>
                        <FormControlLabel  value="service" control={<Radio/>} label="Service" />
                        <FormControlLabel  value="job" control={<Radio/>} label="Job" />
                    </RadioGroup>
                )} />
                
            </FormControl>
            <TextField label="Input Topic" {...register("inputTopic")} />
            <TextField label="Output Topic" {...register("outputTopic")} />
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