import { AppBar, Box, CssBaseline, Toolbar, Typography, Container, Tabs, Tab } from '@mui/material';
import { Link, Route, Routes, useLocation } from 'react-router-dom';
import DiagnosisPage from './pages/DiagnosisPage';
import AppointmentPage from './pages/AppointmentPage';
import CustomerListPage from './pages/CustomerListPage';
import CustomerDetailPage from './pages/CustomerDetailPage';

function a11yTabProps(path: string) {
  return {
    id: `nav-tab-${path}`,
    'aria-controls': `nav-tabpanel-${path}`,
  };
}

const App = () => {
  const location = useLocation();

  const currentTab = (() => {
    if (location.pathname.startsWith('/appointments')) return '/appointments';
    if (location.pathname.startsWith('/customers')) return '/customers';
    if (location.pathname.startsWith('/diagnoses')) return '/diagnoses';
    return '/diagnoses';
  })();

  return (
    <>
      <CssBaseline />
      <AppBar position="static" color="primary">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            TailorCloud – Suit-MBTI Console
          </Typography>
        </Toolbar>
      </AppBar>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={currentTab} variant="fullWidth">
          <Tab
            label="診断一覧"
            value="/diagnoses"
            component={Link}
            to="/diagnoses"
            {...a11yTabProps('diagnoses')}
          />
          <Tab
            label="予約一覧"
            value="/appointments"
            component={Link}
            to="/appointments"
            {...a11yTabProps('appointments')}
          />
          <Tab
            label="顧客一覧"
            value="/customers"
            component={Link}
            to="/customers"
            {...a11yTabProps('customers')}
          />
        </Tabs>
      </Box>
      <Container
        maxWidth="lg"
        sx={{
          mt: { xs: 2, sm: 4 },
          mb: { xs: 2, sm: 4 },
          px: { xs: 1, sm: 2 },
        }}
      >
        <Routes>
          <Route path="/diagnoses" element={<DiagnosisPage />} />
          <Route path="/appointments" element={<AppointmentPage />} />
          <Route path="/customers" element={<CustomerListPage />} />
          <Route path="/customers/:id" element={<CustomerDetailPage />} />
          <Route path="*" element={<DiagnosisPage />} />
        </Routes>
      </Container>
    </>
  );
};

export default App;


