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
import { getDiagnosesByTenant } from '../api/diagnoses';
import type { Diagnosis } from '../types';

const TENANT_ID = import.meta.env.VITE_TENANT_ID || 'tenant_test_suit_mbti';

const DiagnosisPage = () => {
  const [tenantId] = useState<string>(TENANT_ID);

  const {
    data,
    isLoading,
    isError,
    error,
    refetch,
  } = useQuery({
    queryKey: ['diagnoses', tenantId],
    queryFn: () => getDiagnosesByTenant(tenantId, { limit: 50 }),
  });

  useEffect(() => {
    // テナントIDが変わったときの将来拡張用
    void refetch();
  }, [tenantId, refetch]);

  const diagnoses: Diagnosis[] = data?.data ?? [];

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        診断一覧
      </Typography>

      {isLoading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {isError && (
        <Alert severity="error" sx={{ mt: 2 }}>
          診断一覧の取得に失敗しました。
          {error instanceof Error && <Box component="span"> 詳細: {error.message}</Box>}
        </Alert>
      )}

      {!isLoading && !isError && diagnoses.length === 0 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          診断データがまだ登録されていません。
        </Alert>
      )}

      {!isLoading && !isError && diagnoses.length > 0 && (
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
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>診断ID</TableCell>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>ユーザーID</TableCell>
                <TableCell sx={{ minWidth: { xs: 80, sm: 100 } }}>アーキタイプ</TableCell>
                <TableCell sx={{ minWidth: { xs: 80, sm: 100 } }}>プラン</TableCell>
                <TableCell sx={{ minWidth: { xs: 120, sm: 160 } }}>作成日時</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {diagnoses.map((d) => (
                <TableRow key={d.id} hover>
                  <TableCell
                    sx={{
                      wordBreak: 'break-word',
                      maxWidth: { xs: 100, sm: 200 },
                    }}
                  >
                    {d.id}
                  </TableCell>
                  <TableCell
                    sx={{
                      wordBreak: 'break-word',
                      maxWidth: { xs: 100, sm: 200 },
                    }}
                  >
                    {d.user_id}
                  </TableCell>
                  <TableCell>{d.archetype}</TableCell>
                  <TableCell>{d.plan_type ?? '-'}</TableCell>
                  <TableCell>
                    {new Date(d.created_at).toLocaleString('ja-JP', {
                      year: 'numeric',
                      month: '2-digit',
                      day: '2-digit',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
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

export default DiagnosisPage;


