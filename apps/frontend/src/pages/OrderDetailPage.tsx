import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {  MapPin, CreditCard, ChevronLeft, XCircle } from 'lucide-react';
import { useApi } from '../lib/api';
import { useToast } from '../contexts/ToastContext';
import type { Order } from '../lib/api';

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  confirmed: 'bg-blue-100 text-blue-800 border-blue-200',
  processing: 'bg-purple-100 text-purple-800 border-purple-200',
  shipped: 'bg-indigo-100 text-indigo-800 border-indigo-200',
  delivered: 'bg-green-100 text-green-800 border-green-200',
  cancelled: 'bg-red-100 text-red-800 border-red-200',
  refunded: 'bg-gray-100 text-gray-800 border-gray-200',
};

const statusSteps = ['pending', 'confirmed', 'processing', 'shipped', 'delivered'];

export default function OrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const api = useApi();
  const toast = useToast();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);
  const [cancelling, setCancelling] = useState(false);

  useEffect(() => {
    if (id) {
      fetchOrder();
    }
  }, [id]);

  const fetchOrder = async () => {
    if (!id) return;
    
    setLoading(true);
    try {
      const response = await api.getOrderById(id);
      if (response.success && response.data) {
        setOrder(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch order:', error);
      toast.error('Failed to load order details');
      navigate('/orders');
    } finally {
      setLoading(false);
    }
  };

  const handleCancelOrder = async () => {
    if (!order || !confirm('Are you sure you want to cancel this order?')) {
      return;
    }

    setCancelling(true);
    try {
      const response = await api.cancelOrder(order.id);
      if (response.success) {
        toast.success('Order cancelled successfully');
        fetchOrder();
      }
    } catch (error: any) {
      console.error('Failed to cancel order:', error);
      toast.error(error.message || 'Failed to cancel order');
    } finally {
      setCancelling(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading order details...</p>
        </div>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Order not found</h2>
          <button
            onClick={() => navigate('/orders')}
            className="text-green-600 hover:underline"
          >
            Back to Orders
          </button>
        </div>
      </div>
    );
  }

  const currentStatusIndex = statusSteps.indexOf(order.status);
  const canCancel = order.status === 'confirmed';

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        <button
          onClick={() => navigate('/orders')}
          className="flex items-center text-gray-600 hover:text-gray-900 mb-6"
        >
          <ChevronLeft className="w-5 h-5" />
          <span className="ml-1">Back to Orders</span>
        </button>

        {/* Header */}
        <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Order #{order.order_number}</h1>
              <p className="text-gray-600 mt-1">
                Placed on {new Date(order.created_at).toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </p>
            </div>
            <div className="flex flex-col items-end gap-2">
              <span
                className={`px-4 py-2 rounded-lg text-sm font-semibold border ${
                  statusColors[order.status as keyof typeof statusColors] || 'bg-gray-100 text-gray-800'
                }`}
              >
                {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
              </span>
              {canCancel && (
                <button
                  onClick={handleCancelOrder}
                  disabled={cancelling}
                  className="flex items-center gap-2 text-red-600 hover:text-red-800 text-sm disabled:opacity-50"
                >
                  <XCircle className="w-4 h-4" />
                  Cancel Order
                </button>
              )}
            </div>
          </div>

          {/* Order Progress */}
          {order.status !== 'cancelled' && order.status !== 'refunded' && (
            <div className="mt-6">
              <div className="flex items-center justify-between">
                {statusSteps.map((status, index) => (
                  <div key={status} className="flex-1 flex items-center">
                    <div className="flex flex-col items-center flex-1">
                      <div
                        className={`w-10 h-10 rounded-full flex items-center justify-center border-2 ${
                          index <= currentStatusIndex
                            ? 'bg-green-600 border-green-600 text-white'
                            : 'bg-white border-gray-300 text-gray-400'
                        }`}
                      >
                        {index < currentStatusIndex ? '✓' : index + 1}
                      </div>
                      <span className="text-xs mt-2 text-center capitalize">{status}</span>
                    </div>
                    {index < statusSteps.length - 1 && (
                      <div
                        className={`h-1 flex-1 ${
                          index < currentStatusIndex ? 'bg-green-600' : 'bg-gray-300'
                        }`}
                      />
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Order Items */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <h2 className="text-xl font-bold mb-4">Order Items</h2>
              <div className="space-y-4">
                {order.items.map((item) => (
                  <div key={item.id} className="flex gap-4 pb-4 border-b last:border-b-0">
                    <img
                      src={item.product_image_url || 'https://via.placeholder.com/100'}
                      alt={item.product_name}
                      className="w-20 h-20 object-cover rounded"
                    />
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-900">{item.product_name}</h3>
                      <p className="text-sm text-gray-600 mt-1">by {item.shop_name}</p>
                      <p className="text-sm text-gray-600 mt-1">Quantity: {item.quantity}</p>
                      <p className="text-lg font-semibold text-green-600 mt-2">
                        NPR {item.subtotal.toFixed(2)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Order Summary & Details */}
          <div className="lg:col-span-1 space-y-6">
            {/* Order Summary */}
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <h2 className="text-xl font-bold mb-4">Order Summary</h2>
              <div className="space-y-2">
                <div className="flex justify-between text-gray-600">
                  <span>Subtotal</span>
                  <span>NPR {order.subtotal.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Shipping</span>
                  <span>NPR {order.shipping_cost.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Tax</span>
                  <span>NPR {order.tax.toFixed(2)}</span>
                </div>
                {order.discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Discount</span>
                    <span>-NPR {order.discount.toFixed(2)}</span>
                  </div>
                )}
                <div className="flex justify-between text-xl font-bold pt-2 border-t">
                  <span>Total</span>
                  <span className="text-green-600">NPR {order.total.toFixed(2)}</span>
                </div>
              </div>
            </div>

            {/* Shipping Address */}
            {order.shipping_address && (
              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <div className="flex items-center gap-2 mb-3">
                  <MapPin className="w-5 h-5 text-gray-600" />
                  <h2 className="text-lg font-bold">Shipping Address</h2>
                </div>
                <div className="text-gray-700">
                  <p className="font-medium">{order.shipping_address.full_name}</p>
                  <p className="text-sm mt-1">{order.shipping_address.phone}</p>
                  <p className="text-sm mt-2">{order.shipping_address.address_line1}</p>
                  {order.shipping_address.address_line2 && (
                    <p className="text-sm">{order.shipping_address.address_line2}</p>
                  )}
                  <p className="text-sm">
                    {order.shipping_address.city}, {order.shipping_address.state || ''} {order.shipping_address.postal_code || ''}
                  </p>
                  <p className="text-sm">{order.shipping_address.country}</p>
                </div>
              </div>
            )}

            {/* Payment Info */}
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <div className="flex items-center gap-2 mb-3">
                <CreditCard className="w-5 h-5 text-gray-600" />
                <h2 className="text-lg font-bold">Payment</h2>
              </div>
              <div className="text-gray-700">
                <p className="text-sm">
                  <span className="font-medium">Method:</span> {order.payment_method || 'Not specified'}
                </p>
                <p className="text-sm mt-1">
                  <span className="font-medium">Status:</span>{' '}
                  <span
                    className={
                      order.payment_status === 'paid'
                        ? 'text-green-600 font-semibold'
                        : order.payment_status === 'failed'
                        ? 'text-red-600 font-semibold'
                        : 'text-yellow-600 font-semibold'
                    }
                  >
                    {order.payment_status.charAt(0).toUpperCase() + order.payment_status.slice(1)}
                  </span>
                </p>
              </div>
            </div>

            {/* Notes */}
            {order.notes && (
              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <h2 className="text-lg font-bold mb-3">Order Notes</h2>
                <p className="text-gray-700 text-sm">{order.notes}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
