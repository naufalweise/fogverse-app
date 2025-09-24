"use client";

import React from 'react';
import { TextField, Button, Box, Typography, FormControl, FormLabel, RadioGroup, FormControlLabel, Radio } from '@mui/material';
import { Controller, SubmitHandler, useForm } from 'react-hook-form'
import IDeployData from './types/DeployData';


type Props = {}

function DeployForm() {
    const { register, handleSubmit, control} = useForm<IDeployData>({
        defaultValues: {
            type: "service"
        }
    })
    const onSubmit: SubmitHandler<IDeployData> = (data) => {
        console.log(data)
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
            <Button variant="contained" type="submit" sx={{ mt: 2 }}>
                Submit
            </Button>
        </Box>
  );
}

export default DeployForm;