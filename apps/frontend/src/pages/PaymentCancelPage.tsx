import { useNavigate } from 'react-router-dom';
import { XCircle } from 'lucide-react';

export default function PaymentCancelPage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center max-w-md">
        <div className="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-6">
          <XCircle className="w-12 h-12 text-red-500" />
        </div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Payment Cancelled</h1>
        <p className="text-gray-600 mb-8">
          Your payment was not completed. No charges were made. You can try again or choose a different payment method.
        </p>

        <div className="space-y-3">
          <button
            onClick={() => navigate('/checkout')}
            className="w-full bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 font-semibold"
          >
            Try Again
          </button>
          <button
            onClick={() => navigate('/cart')}
            className="w-full border border-gray-300 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-50"
          >
            Back to Cart
          </button>
        </div>
      </div>
    </div>
  );
}
