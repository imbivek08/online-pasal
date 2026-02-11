import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { CheckCircle, Loader2 } from 'lucide-react';
import { useApi } from '../lib/api';
import type { StripeSessionStatus } from '../lib/api';

export default function PaymentSuccessPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const api = useApi();

  const [loading, setLoading] = useState(true);
  const [session, setSession] = useState<StripeSessionStatus | null>(null);
  const [error, setError] = useState<string | null>(null);

  const sessionId = searchParams.get('session_id');

  useEffect(() => {
    if (!sessionId) {
      setError('No session ID found');
      setLoading(false);
      return;
    }

    const verifyPayment = async () => {
      try {
        const response = await api.verifyStripeSession(sessionId);
        if (response.success && response.data) {
          setSession(response.data);
        } else {
          setError('Could not verify payment');
        }
      } catch {
        setError('Failed to verify payment status');
      } finally {
        setLoading(false);
      }
    };

    verifyPayment();
  }, [sessionId]);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-12 h-12 text-green-600 animate-spin mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-gray-900">Verifying your payment...</h2>
          <p className="text-gray-600 mt-2">Please wait while we confirm your order.</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center max-w-md">
          <div className="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">⚠️</span>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Payment Verification Issue</h2>
          <p className="text-gray-600 mb-6">{error}</p>
          <p className="text-sm text-gray-500 mb-6">
            Don't worry — if your payment went through, your order will be updated automatically.
            Check your orders page for the latest status.
          </p>
          <button
            onClick={() => navigate('/orders')}
            className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700"
          >
            View My Orders
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center max-w-md">
        <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
          <CheckCircle className="w-12 h-12 text-green-600" />
        </div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Payment Successful!</h1>
        <p className="text-gray-600 mb-2">Thank you for your purchase.</p>
        {session && (
          <p className="text-sm text-gray-500 mb-8">
            Order <span className="font-semibold text-gray-700">#{session.order_number}</span> has been confirmed.
          </p>
        )}

        <div className="space-y-3">
          {session && (
            <button
              onClick={() => navigate(`/orders/${session.order_id}`)}
              className="w-full bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 font-semibold"
            >
              View Order Details
            </button>
          )}
          <button
            onClick={() => navigate('/products')}
            className="w-full border border-gray-300 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-50"
          >
            Continue Shopping
          </button>
        </div>
      </div>
    </div>
  );
}
