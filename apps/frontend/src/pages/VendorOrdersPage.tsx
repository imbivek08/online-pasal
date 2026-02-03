import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Clock, CheckCircle, XCircle } from 'lucide-react';
import { useApi } from '../lib/api';
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

const statusIcons = {
  pending: Clock,
  confirmed: CheckCircle,
  processing: Package,
  shipped: Package,
  delivered: CheckCircle,
  cancelled: XCircle,
  refunded: XCircle,
};

export default function VendorOrdersPage() {
  const navigate = useNavigate();
  const api = useApi();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [updatingStatus, setUpdatingStatus] = useState<string | null>(null);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const response = await api.getVendorOrders();
      if (response.success && response.data) {
        setOrders(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch vendor orders:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleStatusUpdate = async (orderId: string, newStatus: string) => {
    setUpdatingStatus(orderId);
    try {
      const response = await api.updateOrderStatus(orderId, newStatus);
      if (response.success) {
        // Update local state
        setOrders(orders.map(order => 
          order.id === orderId 
            ? { ...order, status: newStatus }
            : order
        ));
      }
    } catch (error: any) {
      console.error('Failed to update order status:', error);
      alert(error.message || 'Failed to update order status');
    } finally {
      setUpdatingStatus(null);
    }
  };

  const getNextStatuses = (currentStatus: string): string[] => {
    const statusFlow: { [key: string]: string[] } = {
      pending: ['confirmed', 'cancelled'],
      confirmed: ['processing', 'cancelled'],
      processing: ['shipped', 'cancelled'],
      shipped: ['delivered'],
      delivered: ['refunded'],
      cancelled: [],
      refunded: [],
    };
    return statusFlow[currentStatus] || [];
  };

  const filteredOrders = statusFilter === 'all' 
    ? orders 
    : orders.filter(order => order.status === statusFilter);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading vendor orders...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Vendor Orders</h1>
          <p className="text-gray-600 mt-2">Manage orders for your shop</p>
        </div>

        {/* Filter Tabs */}
        <div className="bg-white rounded-lg border border-gray-200 p-4 mb-6">
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setStatusFilter('all')}
              className={`px-4 py-2 rounded-lg font-medium transition ${
                statusFilter === 'all'
                  ? 'bg-green-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              All ({orders.length})
            </button>
            {Object.keys(statusColors).map((status) => (
              <button
                key={status}
                onClick={() => setStatusFilter(status)}
                className={`px-4 py-2 rounded-lg font-medium transition capitalize ${
                  statusFilter === status
                    ? 'bg-green-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {status} ({orders.filter(o => o.status === status).length})
              </button>
            ))}
          </div>
        </div>

        {/* Orders List */}
        {filteredOrders.length === 0 ? (
          <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
            <Package className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-2xl font-bold text-gray-900 mb-2">No orders found</h2>
            <p className="text-gray-600">
              {statusFilter === 'all'
                ? "You don't have any orders yet."
                : `No ${statusFilter} orders found.`}
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredOrders.map((order) => {
              const StatusIcon = statusIcons[order.status as keyof typeof statusIcons] || Package;
              const nextStatuses = getNextStatuses(order.status);

              return (
                <div
                  key={order.id}
                  className="bg-white rounded-lg border border-gray-200 p-6 hover:shadow-md transition"
                >
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <h3 className="text-lg font-bold text-gray-900">
                          Order #{order.order_number}
                        </h3>
                        <span
                          className={`px-3 py-1 rounded-lg text-xs font-semibold border capitalize flex items-center gap-1 ${
                            statusColors[order.status as keyof typeof statusColors]
                          }`}
                        >
                          <StatusIcon className="w-3 h-3" />
                          {order.status}
                        </span>
                      </div>
                      <p className="text-sm text-gray-600">
                        Placed {new Date(order.created_at).toLocaleDateString('en-US', {
                          year: 'numeric',
                          month: 'short',
                          day: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </p>
                      <p className="text-sm text-gray-600 mt-1">
                        Payment: {order.payment_method || 'N/A'} • Status:{' '}
                        <span
                          className={
                            order.payment_status === 'paid'
                              ? 'text-green-600 font-semibold'
                              : order.payment_status === 'failed'
                              ? 'text-red-600 font-semibold'
                              : 'text-yellow-600 font-semibold'
                          }
                        >
                          {order.payment_status}
                        </span>
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-2xl font-bold text-green-600">
                        NPR {order.total.toFixed(2)}
                      </p>
                      <p className="text-sm text-gray-600 mt-1">
                        {order.items.reduce((sum, item) => sum + item.quantity, 0)} items
                      </p>
                    </div>
                  </div>

                  {/* Order Items */}
                  <div className="flex gap-2 mb-4 overflow-x-auto">
                    {order.items.slice(0, 5).map((item) => (
                      <img
                        key={item.id}
                        src={item.product_image_url || 'https://via.placeholder.com/80'}
                        alt={item.product_name}
                        className="w-20 h-20 object-cover rounded border border-gray-200"
                      />
                    ))}
                    {order.items.length > 5 && (
                      <div className="w-20 h-20 bg-gray-100 rounded border border-gray-200 flex items-center justify-center">
                        <span className="text-gray-600 text-sm font-semibold">
                          +{order.items.length - 5}
                        </span>
                      </div>
                    )}
                  </div>

                  {/* Shipping Address */}
                  {order.shipping_address && (
                    <div className="bg-gray-50 rounded p-3 mb-4">
                      <p className="text-xs text-gray-600 font-semibold mb-1">SHIPPING TO:</p>
                      <p className="text-sm text-gray-900 font-medium">
                        {order.shipping_address.full_name}
                      </p>
                      <p className="text-sm text-gray-700">
                        {order.shipping_address.address_line1}
                        {order.shipping_address.address_line2 && `, ${order.shipping_address.address_line2}`}
                      </p>
                      <p className="text-sm text-gray-700">
                        {order.shipping_address.city}, {order.shipping_address.state} {order.shipping_address.postal_code}
                      </p>
                      <p className="text-sm text-gray-700">{order.shipping_address.phone}</p>
                    </div>
                  )}

                  {/* Actions */}
                  <div className="flex items-center justify-between pt-4 border-t">
                    <button
                      onClick={() => navigate(`/orders/${order.id}`)}
                      className="text-green-600 hover:text-green-700 font-medium text-sm"
                    >
                      View Details →
                    </button>

                    {nextStatuses.length > 0 && (
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-gray-600">Update status:</span>
                        {nextStatuses.map((status) => (
                          <button
                            key={status}
                            onClick={() => handleStatusUpdate(order.id, status)}
                            disabled={updatingStatus === order.id}
                            className={`px-3 py-1 rounded text-xs font-semibold capitalize transition disabled:opacity-50 ${
                              status === 'cancelled'
                                ? 'bg-red-100 text-red-700 hover:bg-red-200 border border-red-200'
                                : 'bg-blue-100 text-blue-700 hover:bg-blue-200 border border-blue-200'
                            }`}
                          >
                            {updatingStatus === order.id ? 'Updating...' : `Mark ${status}`}
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
