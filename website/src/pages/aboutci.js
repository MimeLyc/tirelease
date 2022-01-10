import * as React from 'react';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';

import Layout from '../layout/Layout';
import CIData from '../components/aboutci/CIData';

const AboutCI = () => {
    return (
        <Layout>
            <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                    <CIData />
                </Paper>
            </Container>
        </Layout>
    )
};

export default AboutCI;