"use client";

import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';

export default function BackButton() {
    const router = useRouter()
    return (<Button startIcon={<ArrowBackIcon />} onClick={() => router.back()}>Back</Button>)
}