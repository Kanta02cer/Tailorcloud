import React from 'react';
import { useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
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
import { getCustomer, getCustomerOrders } from '../api/customers';
import { createOrder, generateOrderDocument } from '../api/orders';
import type { Customer, OrderSummary } from '../types/customer';
import type { CreateOrderRequest } from '../types/order';

const formatDateInput = (date: Date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const CustomerDetailPage = () => {
  const { id } = useParams<{ id: string }>();

  const {
    data: customer,
    isLoading: isCustomerLoading,
    isError: isCustomerError,
    error: customerError,
  } = useQuery({
    queryKey: ['customer', id],
    enabled: Boolean(id),
    queryFn: () => getCustomer(id as string),
  });

  const {
    data: orders,
    isLoading: isOrdersLoading,
    isError: isOrdersError,
    error: ordersError,
  } = useQuery({
    queryKey: ['customerOrders', id],
    enabled: Boolean(id),
    queryFn: () => getCustomerOrders(id as string),
  });

  if (!id) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        顧客IDが指定されていません。
      </Alert>
    );
  }

  if (isCustomerLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (isCustomerError) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        顧客情報の取得に失敗しました。
        {customerError instanceof Error && (
          <Box component="span"> 詳細: {customerError.message}</Box>
        )}
      </Alert>
    );
  }

  const typedCustomer = customer as Customer | undefined;
  const typedOrders = (orders ?? []) as OrderSummary[];
  const queryClient = useQueryClient();

  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [fabricId, setFabricId] = React.useState('');
  const [totalAmount, setTotalAmount] = React.useState<string>('');
  const [deliveryDate, setDeliveryDate] = React.useState<string>(formatDateInput(new Date()));
  const [description, setDescription] = React.useState('オーダースーツ縫製');
  const [lastGeneratedDocUrl, setLastGeneratedDocUrl] = React.useState<string | null>(null);
  const [validationError, setValidationError] = React.useState<string | null>(null);

  const createOrderMutation = useMutation({
    mutationFn: (req: CreateOrderRequest) => createOrder(req),
    onSuccess: async (order) => {
      // 注文作成成功後、注文履歴を再取得
      await queryClient.invalidateQueries({ queryKey: ['customerOrders', id] });
      // PDF生成を自動実行
      try {
        const doc = await generateDocumentMutation.mutateAsync(order.id);
        setLastGeneratedDocUrl(doc.doc_url);
        setIsDialogOpen(false);
        setValidationError(null);
      } catch (error) {
        // PDF生成エラーは個別に表示
        console.error('PDF生成エラー:', error);
      }
    },
    onError: (error) => {
      console.error('発注作成エラー:', error);
      setValidationError(error instanceof Error ? error.message : '発注作成に失敗しました');
    },
  });

  const generateDocumentMutation = useMutation({
    mutationFn: (orderId: string) => generateOrderDocument(orderId),
    onError: (error) => {
      console.error('PDF生成エラー:', error);
      setValidationError(error instanceof Error ? error.message : 'PDF生成に失敗しました');
    },
  });

  if (!typedCustomer) {
    return (
      <Alert severity="info" sx={{ mt: 2 }}>
        顧客が見つかりませんでした。
      </Alert>
    );
  }

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        顧客詳細
      </Typography>

      <Paper sx={{ p: { xs: 1.5, sm: 2 }, mb: 3 }}>
        <Typography variant="h6" sx={{ fontSize: { xs: '1.1rem', sm: '1.25rem' } }}>
          {typedCustomer.name}
        </Typography>
        <Box sx={{ mt: 1, display: 'flex', flexDirection: 'column', gap: 0.5 }}>
          <Typography variant="body2" sx={{ fontSize: { xs: '0.85rem', sm: '0.875rem' } }}>
            メール: {typedCustomer.email ?? '-'}
          </Typography>
          <Typography variant="body2" sx={{ fontSize: { xs: '0.85rem', sm: '0.875rem' } }}>
            電話番号: {typedCustomer.phone ?? '-'}
          </Typography>
          <Typography variant="body2" sx={{ fontSize: { xs: '0.85rem', sm: '0.875rem' } }}>
            住所: {typedCustomer.address ?? '-'}
          </Typography>
          <Typography variant="body2" sx={{ fontSize: { xs: '0.85rem', sm: '0.875rem' } }}>
            登録日: {new Date(typedCustomer.created_at).toLocaleDateString()}
          </Typography>
        </Box>

        <Box
          sx={{
            mt: 2,
            display: 'flex',
            flexDirection: { xs: 'column', sm: 'row' },
            gap: 2,
            alignItems: { xs: 'stretch', sm: 'center' },
          }}
        >
          <Button
            variant="contained"
            onClick={() => {
              setIsDialogOpen(true);
              setValidationError(null);
            }}
            fullWidth={false}
            sx={{ minWidth: { xs: '100%', sm: 'auto' } }}
          >
            この顧客で発注を作成
          </Button>
          {lastGeneratedDocUrl && (
            <Button
              variant="outlined"
              color="secondary"
              href={lastGeneratedDocUrl}
              target="_blank"
              rel="noopener noreferrer"
              fullWidth={false}
              sx={{ minWidth: { xs: '100%', sm: 'auto' } }}
            >
              最新の発注書PDFを開く
            </Button>
          )}
        </Box>
      </Paper>

      <Divider sx={{ mb: 2 }} />

      <Typography variant="h6" gutterBottom>
        注文履歴
      </Typography>

      {isOrdersLoading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
          <CircularProgress size={24} />
        </Box>
      )}

      {isOrdersError && (
        <Alert severity="error" sx={{ mt: 2 }}>
          注文履歴の取得に失敗しました。
          {ordersError instanceof Error && <Box component="span"> 詳細: {ordersError.message}</Box>}
        </Alert>
      )}

      {!isOrdersLoading && !isOrdersError && typedOrders.length === 0 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          この顧客の注文履歴はまだありません。
        </Alert>
      )}

      {!isOrdersLoading && !isOrdersError && typedOrders.length > 0 && (
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
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>注文ID</TableCell>
                <TableCell sx={{ minWidth: { xs: 80, sm: 100 } }}>ステータス</TableCell>
                <TableCell align="right" sx={{ minWidth: { xs: 100, sm: 120 } }}>
                  金額
                </TableCell>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>納期</TableCell>
                <TableCell sx={{ minWidth: { xs: 100, sm: 120 } }}>作成日</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {typedOrders.map((o) => (
                <TableRow key={o.id} hover>
                  <TableCell
                    sx={{
                      wordBreak: 'break-word',
                      maxWidth: { xs: 100, sm: 200 },
                    }}
                  >
                    {o.id}
                  </TableCell>
                  <TableCell>{o.status}</TableCell>
                  <TableCell align="right">
                    {o.total_amount.toLocaleString()} 円
                  </TableCell>
                  <TableCell>
                    {o.delivery_date ? new Date(o.delivery_date).toLocaleDateString() : '-'}
                  </TableCell>
                  <TableCell>{new Date(o.created_at).toLocaleDateString()}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      <Dialog
        open={isDialogOpen}
        onClose={() => {
          if (!createOrderMutation.isPending && !generateDocumentMutation.isPending) {
            setIsDialogOpen(false);
            setValidationError(null);
          }
        }}
        fullWidth
        maxWidth="sm"
        PaperProps={{
          sx: {
            m: { xs: 1, sm: 2 },
            maxHeight: { xs: '90vh', sm: 'auto' },
          },
        }}
      >
        <DialogTitle sx={{ fontSize: { xs: '1rem', sm: '1.25rem' } }}>
          新規発注の作成
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
              顧客: {typedCustomer.name}
            </Typography>
            <TextField
              label="生地ID"
              value={fabricId}
              onChange={(e) => {
                setFabricId(e.target.value);
                setValidationError(null);
              }}
              fullWidth
              size="small"
              required
              error={!fabricId && validationError !== null}
              helperText={!fabricId ? '生地IDは必須です' : ''}
              inputProps={{
                'aria-label': '生地ID',
              }}
            />
            <TextField
              label="金額（円）"
              type="number"
              value={totalAmount}
              onChange={(e) => {
                setTotalAmount(e.target.value);
                setValidationError(null);
              }}
              fullWidth
              size="small"
              required
              error={(!totalAmount || Number(totalAmount) <= 0) && validationError !== null}
              helperText={
                !totalAmount || Number(totalAmount) <= 0 ? '金額は1円以上を入力してください' : ''
              }
              inputProps={{
                min: 1,
                step: 1,
                'aria-label': '金額',
              }}
            />
            <TextField
              label="納期"
              type="date"
              value={deliveryDate}
              onChange={(e) => {
                setDeliveryDate(e.target.value);
                setValidationError(null);
              }}
              fullWidth
              size="small"
              required
              InputLabelProps={{ shrink: true }}
              inputProps={{
                'aria-label': '納期',
              }}
            />
            <TextField
              label="内容（給付の内容）"
              value={description}
              onChange={(e) => {
                setDescription(e.target.value);
                setValidationError(null);
              }}
              fullWidth
              size="small"
              required
              multiline
              minRows={2}
              inputProps={{
                'aria-label': '給付の内容',
              }}
            />
          </Box>
          {validationError && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {validationError}
            </Alert>
          )}
          {createOrderMutation.isError && !validationError && (
            <Alert severity="error" sx={{ mt: 2 }}>
              発注作成に失敗しました。
              {createOrderMutation.error instanceof Error && (
                <Box component="span" sx={{ display: 'block', mt: 0.5, fontSize: '0.875rem' }}>
                  詳細: {createOrderMutation.error.message}
                </Box>
              )}
            </Alert>
          )}
          {generateDocumentMutation.isError && !validationError && (
            <Alert severity="error" sx={{ mt: 2 }}>
              発注書PDFの生成に失敗しました。
              {generateDocumentMutation.error instanceof Error && (
                <Box component="span" sx={{ display: 'block', mt: 0.5, fontSize: '0.875rem' }}>
                  詳細: {generateDocumentMutation.error.message}
                </Box>
              )}
            </Alert>
          )}
        </DialogContent>
        <DialogActions
          sx={{
            flexDirection: { xs: 'column-reverse', sm: 'row' },
            gap: { xs: 1, sm: 0 },
            px: { xs: 2, sm: 3 },
            pb: { xs: 2, sm: 2 },
          }}
        >
          <Button
            onClick={() => {
              setIsDialogOpen(false);
              setValidationError(null);
            }}
            disabled={createOrderMutation.isPending || generateDocumentMutation.isPending}
            fullWidth={false}
            sx={{ minWidth: { xs: '100%', sm: 'auto' } }}
          >
            キャンセル
          </Button>
          <Button
            variant="contained"
            onClick={async () => {
              if (!typedCustomer) return;

              // バリデーション
              if (!fabricId.trim()) {
                setValidationError('生地IDは必須です');
                return;
              }
              const amount = Number(totalAmount);
              if (!Number.isFinite(amount) || amount <= 0) {
                setValidationError('金額は1円以上の数値を入力してください');
                return;
              }
              if (!deliveryDate) {
                setValidationError('納期を選択してください');
                return;
              }
              if (!description.trim()) {
                setValidationError('給付の内容を入力してください');
                return;
              }

              setValidationError(null);

              // ISO 8601 (RFC3339) 形式に変換: "2025-12-31T00:00:00Z"
              const dateObj = new Date(deliveryDate);
              if (isNaN(dateObj.getTime())) {
                setValidationError('納期の形式が不正です');
                return;
              }
              const isoDeliveryDate = dateObj.toISOString();

              const req: CreateOrderRequest = {
                customer_id: typedCustomer.id,
                fabric_id: fabricId.trim(),
                total_amount: amount,
                delivery_date: isoDeliveryDate,
                details: {
                  description: description.trim() || 'オーダースーツ縫製',
                  measurement_data: {},
                  adjustments: {},
                },
              };

              // mutationはonSuccessでPDF生成まで自動実行される
              await createOrderMutation.mutateAsync(req);
            }}
            disabled={createOrderMutation.isPending || generateDocumentMutation.isPending}
            fullWidth={false}
            sx={{ minWidth: { xs: '100%', sm: 'auto' } }}
          >
            {createOrderMutation.isPending || generateDocumentMutation.isPending
              ? '処理中...'
              : '発注作成 & PDF生成'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CustomerDetailPage;


