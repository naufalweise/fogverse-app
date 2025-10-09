import * as React from 'react';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import DeployForm from './DeployForm';
import BackButton from '@/components/BackButton';


export default function Deploy() {
  return (
    <Container maxWidth="lg">
      <BackButton /> 
      <Box
        sx={{
          my: 4,
          display: 'flex',
          flexDirection: 'column',
          gap: '1em', 
        }}
      >
        <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
        Create Job
        </Typography>
        <DeployForm />
      </Box>
    </Container>
  );
}