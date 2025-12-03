import { useEffect, useState } from 'react';
import {
  Alert,
  Box,
  CircularProgress,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { listAppointments } from '../api/appointments';
import type { Appointment } from '../types';

const TENANT_ID = import.meta.env.VITE_TENANT_ID || 'tenant_test_suit_mbti';

const AppointmentPage = () => {
  const [tenantId] = useState<string>(TENANT_ID);

  const {
    data,
    isLoading,
    isError,
    error,
    refetch,
  } = useQuery({
    queryKey: ['appointments', tenantId],
    queryFn: () => listAppointments(tenantId),
  });

  useEffect(() => {
    void refetch();
  }, [tenantId, refetch]);

  const appointments: Appointment[] = data?.data ?? [];

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        予約一覧
      </Typography>

      {isLoading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {isError && (
        <Alert severity="error" sx={{ mt: 2 }}>
          予約一覧の取得に失敗しました。
          {error instanceof Error && <Box component="span"> 詳細: {error.message}</Box>}
        </Alert>
      )}

      {!isLoading && !isError && appointments.length === 0 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          予約データがまだ登録されていません。
        </Alert>
      )}

      {!isLoading && !isError && appointments.length > 0 && (
        <TableContainer
          component={Paper}
          sx={{
            mt: 2,
            overflowX: 'auto',
            '& .MuiTableCell-root': {
              fontSize: { xs: '0.75rem', sm: '0.875rem' },
              padding: { xs: '8px 4px', sm: '16px' },
            },
          }}
        >
          <Table size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>予約ID</TableCell>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 }, display: { xs: 'none', md: 'table-cell' } }}>
                  ユーザーID
                </TableCell>
                <TableCell sx={{ minWidth: { xs: 80, sm: 100 }, display: { xs: 'none', sm: 'table-cell' } }}>
                  フィッターID
                </TableCell>
                <TableCell sx={{ minWidth: { xs: 140, sm: 180 } }}>開始日時</TableCell>
                <TableCell sx={{ minWidth: { xs: 70, sm: 90 } }} align="right">
                  時間(分)
                </TableCell>
                <TableCell sx={{ minWidth: { xs: 80, sm: 100 } }}>ステータス</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {appointments.map((a) => (
                <TableRow key={a.id} hover>
                  <TableCell
                    sx={{
                      wordBreak: 'break-word',
                      maxWidth: { xs: 100, sm: 200 },
                    }}
                  >
                    {a.id}
                  </TableCell>
                  <TableCell
                    sx={{
                      wordBreak: 'break-word',
                      display: { xs: 'none', md: 'table-cell' },
                    }}
                  >
                    {a.user_id}
                  </TableCell>
                  <TableCell
                    sx={{
                      display: { xs: 'none', sm: 'table-cell' },
                    }}
                  >
                    {a.fitter_id ?? '-'}
                  </TableCell>
                  <TableCell>
                    {new Date(a.appointment_datetime).toLocaleString('ja-JP', {
                      year: 'numeric',
                      month: '2-digit',
                      day: '2-digit',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </TableCell>
                  <TableCell align="right">{a.duration_minutes}</TableCell>
                  <TableCell>
                    <Box
                      component="span"
                      sx={{
                        px: 1,
                        py: 0.5,
                        borderRadius: 1,
                        fontSize: '0.75rem',
                        backgroundColor:
                          a.status === 'Confirmed'
                            ? 'success.light'
                            : a.status === 'Cancelled'
                              ? 'error.light'
                              : a.status === 'Completed'
                                ? 'info.light'
                                : 'grey.300',
                        color: 'text.primary',
                      }}
                    >
                      {a.status}
                    </Box>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Box>
  );
};

export default AppointmentPage;


