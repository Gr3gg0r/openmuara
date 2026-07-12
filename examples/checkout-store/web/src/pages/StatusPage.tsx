import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { CheckCircle, XCircle, Loader2, ShoppingBag } from 'lucide-react';
import { fetchPayment } from '../api/client';
import type { Payment } from '../types';

function formatPrice(price: number, currency: string) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(price);
}

export default function StatusPage({ variant }: { variant: 'success' | 'cancel' }) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const ref = searchParams.get('ref');
  const [payment, setPayment] = useState<Payment | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!ref) {
      setError('No payment reference provided');
      setLoading(false);
      return;
    }

    fetchPayment(ref)
      .then(setPayment)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [ref]);

  return (
    <div className="min-h-screen bg-base-200">
      <div className="navbar bg-base-100 shadow-sm">
        <div className="container mx-auto px-4">
          <div className="flex-1">
            <button
              className="btn btn-ghost text-xl font-bold"
              onClick={() => navigate('/')}
            >
              <ShoppingBag className="h-6 w-6 text-primary" />
              OpenMuara Store
            </button>
          </div>
        </div>
      </div>

      <div className="container mx-auto flex items-center justify-center px-4 py-16">
        <div className="card w-full max-w-lg bg-base-100 shadow-xl">
          <div className="card-body items-center text-center">
            {loading ? (
              <>
                <Loader2 className="h-16 w-16 animate-spin text-primary" />
                <h2 className="mt-4 text-2xl font-bold">Checking payment status...</h2>
              </>
            ) : error ? (
              <>
                <XCircle className="h-16 w-16 text-error" />
                <h2 className="mt-4 text-2xl font-bold">Something went wrong</h2>
                <p className="text-base-content/70">{error}</p>
              </>
            ) : variant === 'success' && payment?.status === 'paid' ? (
              <>
                <div className="rounded-full bg-success/10 p-4">
                  <CheckCircle className="h-16 w-16 text-success" />
                </div>
                <h2 className="mt-4 text-3xl font-bold">Payment successful!</h2>
                <p className="mt-2 text-base-content/70">
                  Thank you, {payment.name}. A confirmation email has been sent to{' '}
                  {payment.email}.
                </p>
                <div className="mt-6 w-full rounded-xl bg-base-200 p-4 text-left">
                  <div className="flex justify-between">
                    <span className="text-base-content/70">Reference</span>
                    <span className="font-mono">{payment.ref}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">Method</span>
                    <span className="capitalize">{payment.method}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">Amount</span>
                    <span className="font-bold">
                      {formatPrice(payment.amount, payment.currency)}
                    </span>
                  </div>
                </div>
              </>
            ) : (
              <>
                <div className="rounded-full bg-error/10 p-4">
                  <XCircle className="h-16 w-16 text-error" />
                </div>
                <h2 className="mt-4 text-3xl font-bold">
                  Payment {payment?.status ?? 'cancelled'}
                </h2>
                <p className="mt-2 text-base-content/70">
                  {ref ? `Reference: ${ref}` : 'Your payment could not be completed.'}
                </p>
              </>
            )}

            <button
              className="btn btn-primary mt-8"
              onClick={() => navigate('/')}
              disabled={loading}
            >
              Back to store
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
