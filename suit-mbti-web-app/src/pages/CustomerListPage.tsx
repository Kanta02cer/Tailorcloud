import { useState } from 'react';
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { listCustomers } from '../api/customers';
import type { Customer } from '../types/customer';

const CustomerListPage = () => {
  const [search, setSearch] = useState('');
  const navigate = useNavigate();

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ['customers', search],
    queryFn: () => listCustomers({ search: search || undefined }),
  });

  const customers: Customer[] = data?.customers ?? [];

  const handleSearch = () => {
    void refetch();
  };

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        顧客一覧
      </Typography>

      <Box
        sx={{
          display: 'flex',
          flexDirection: { xs: 'column', sm: 'row' },
          gap: 2,
          mb: 2,
        }}
      >
        <TextField
          size="small"
          label="名前・電話で検索"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              handleSearch();
            }
          }}
          fullWidth={false}
          sx={{ flex: { xs: 1, sm: '0 1 auto' }, minWidth: { xs: '100%', sm: 200 } }}
        />
        <Button
          variant="contained"
          onClick={handleSearch}
          sx={{ minWidth: { xs: '100%', sm: 'auto' } }}
        >
          検索
        </Button>
      </Box>

      {isLoading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {isError && (
        <Alert severity="error" sx={{ mt: 2 }}>
          顧客一覧の取得に失敗しました。
          {error instanceof Error && <Box component="span"> 詳細: {error.message}</Box>}
        </Alert>
      )}

      {!isLoading && !isError && customers.length === 0 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          顧客データがまだ登録されていません。
        </Alert>
      )}

      {!isLoading && !isError && customers.length > 0 && (
        <TableContainer
          component={Paper}
          sx={{
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
                <TableCell sx={{ minWidth: { xs: 120, sm: 150 } }}>顧客名</TableCell>
                <TableCell sx={{ minWidth: { xs: 150, sm: 200 }, display: { xs: 'none', md: 'table-cell' } }}>
                  メールアドレス
                </TableCell>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>電話番号</TableCell>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>登録日</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {customers.map((c) => (
                <TableRow
                  key={c.id}
                  hover
                  sx={{
                    cursor: 'pointer',
                    '&:active': {
                      backgroundColor: 'action.selected',
                    },
                  }}
                  onClick={() => navigate(`/customers/${c.id}`)}
                >
                  <TableCell
                    sx={{
                      fontWeight: 'medium',
                      wordBreak: 'break-word',
                    }}
                  >
                    {c.name}
                  </TableCell>
                  <TableCell
                    sx={{
                      wordBreak: 'break-word',
                      display: { xs: 'none', md: 'table-cell' },
                    }}
                  >
                    {c.email ?? '-'}
                  </TableCell>
                  <TableCell>{c.phone ?? '-'}</TableCell>
                  <TableCell>
                    {new Date(c.created_at).toLocaleDateString('ja-JP')}
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

export default CustomerListPage;


