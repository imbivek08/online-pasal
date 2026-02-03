import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, ChevronRight } from 'lucide-react';
import { useApi } from '../lib/api';
import type { Order } from '../lib/api';

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  confirmed: 'bg-blue-100 text-blue-800',
  processing: 'bg-purple-100 text-purple-800',
  shipped: 'bg-indigo-100 text-indigo-800',
  delivered: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
  refunded: 'bg-gray-100 text-gray-800',
};

export default function OrdersPage() {
  const navigate = useNavigate();
  const api = useApi();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const response = await api.getOrders();
      if (response.success && response.data) {
        setOrders(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch orders:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading your orders...</p>
        </div>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center max-w-md px-4">
          <Package className="w-24 h-24 text-gray-300 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-2">No orders yet</h2>
          <p className="text-gray-600 mb-6">
            You haven't placed any orders. Start shopping to see your orders here.
          </p>
          <button
            onClick={() => navigate('/products')}
            className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 transition-colors"
          >
            Start Shopping
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">My Orders</h1>

        <div className="space-y-4">
          {orders.map((order) => (
            <div
              key={order.id}
              onClick={() => navigate(`/orders/${order.id}`)}
              className="bg-white rounded-lg border border-gray-200 p-6 hover:shadow-md transition-shadow cursor-pointer"
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">
                    Order #{order.order_number}
                  </h3>
                  <p className="text-sm text-gray-500 mt-1">
                    Placed on {new Date(order.created_at).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric',
                    })}
                  </p>
                </div>
                <div className="flex items-center gap-3">
                  <span
                    className={`px-3 py-1 rounded-full text-sm font-medium ${
                      statusColors[order.status as keyof typeof statusColors] || 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                  </span>
                  <ChevronRight className="w-5 h-5 text-gray-400" />
                </div>
              </div>

              <div className="flex items-center gap-4 mb-4">
                {order.items && order.items.slice(0, 3).map((item, index) => (
                  <img
                    key={index}
                    src={item.product_image_url || 'https://via.placeholder.com/80'}
                    alt={item.product_name}
                    className="w-16 h-16 object-cover rounded border border-gray-200"
                  />
                ))}
                {order.items && order.items.length > 3 && (
                  <div className="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-gray-600 font-medium">
                    +{order.items.length - 3}
                  </div>
                )}
              </div>

              <div className="flex items-center justify-between pt-4 border-t border-gray-200">
                <div>
                  <p className="text-sm text-gray-600">
                    {order.items?.length || 0} {order.items?.length === 1 ? 'item' : 'items'}
                  </p>
                  <p className="text-sm text-gray-600 mt-1">
                    Payment: {order.payment_method || 'Not specified'}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-600">Total</p>
                  <p className="text-2xl font-bold text-green-600">NPR {order.total.toFixed(2)}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
