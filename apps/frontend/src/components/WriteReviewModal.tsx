import { useState, useEffect } from 'react';
import { X } from 'lucide-react';
import RatingStars from './RatingStars';
import type { CreateReviewRequest, Order } from '../lib/api';
import { useApi } from '../lib/api';

interface WriteReviewModalProps {
  productId: string;
  productName: string;
  orderId: string;
  onClose: () => void;
  onSubmit: (review: CreateReviewRequest) => Promise<void>;
}

export default function WriteReviewModal({
  productId,
  productName,
  orderId: initialOrderId,
  onClose,
  onSubmit
}: WriteReviewModalProps) {
  const api = useApi();
  const [rating, setRating] = useState(5);
  const [title, setTitle] = useState('');
  const [comment, setComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [orders, setOrders] = useState<Order[]>([]);
  const [selectedOrderId, setSelectedOrderId] = useState(initialOrderId);
  const [loadingOrders, setLoadingOrders] = useState(true);

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const response = await api.getOrders();
        if (response.success && response.data) {
          // Filter for delivered orders that contain this product
          const deliveredOrders = response.data.filter(
            (order: Order) => 
              order.status === 'delivered' && 
              order.items.some(item => item.product_id === productId)
          );
          setOrders(deliveredOrders);
          
          // Auto-select first order if not pre-selected
          if (!initialOrderId && deliveredOrders.length > 0) {
            setSelectedOrderId(deliveredOrders[0].id);
          }
        }
      } catch (error) {
        console.error('Failed to fetch orders:', error);
      } finally {
        setLoadingOrders(false);
      }
    };

    fetchOrders();
  }, [productId, initialOrderId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!selectedOrderId) {
      alert('Please select an order to review');
      return;
    }

    setIsSubmitting(true);

    try {
      await onSubmit({
        product_id: productId,
        order_id: selectedOrderId,
        rating,
        title: title || undefined,
        comment: comment || undefined,
      });
      onClose();
    } catch (error) {
      console.error('Failed to submit review:', error);
      alert('Failed to submit review. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Write a Review</h2>
            <p className="text-gray-600 mt-1">{productName}</p>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {/* Order Selection */}
          {loadingOrders ? (
            <div className="text-center py-4">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600 mx-auto"></div>
              <p className="text-gray-600 mt-2">Loading your orders...</p>
            </div>
          ) : orders.length > 0 ? (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Select Order *
              </label>
              <select
                value={selectedOrderId}
                onChange={(e) => setSelectedOrderId(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                required
              >
                <option value="">Choose an order</option>
                {orders.map((order) => (
                  <option key={order.id} value={order.id}>
                    Order #{order.order_number} - {new Date(order.delivered_at || order.created_at).toLocaleDateString()}
                  </option>
                ))}
              </select>
              <p className="text-sm text-gray-500 mt-1">
                Only delivered orders containing this product are shown
              </p>
            </div>
          ) : (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <p className="text-yellow-800">
                You need to purchase and receive this product before you can review it.
              </p>
            </div>
          )}

          {/* Rating */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Rating *
            </label>
            <RatingStars
              rating={rating}
              size="lg"
              interactive
              onRatingChange={setRating}
            />
            <p className="text-sm text-gray-500 mt-1">
              Click to rate from 1 to 5 stars
            </p>
          </div>

          {/* Title */}
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
              Review Title
            </label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={200}
              placeholder="Sum up your experience in one line"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            />
            <p className="text-sm text-gray-500 mt-1">
              {title.length}/200 characters
            </p>
          </div>

          {/* Comment */}
          <div>
            <label htmlFor="comment" className="block text-sm font-medium text-gray-700 mb-2">
              Review
            </label>
            <textarea
              id="comment"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              maxLength={2000}
              rows={6}
              placeholder="Share your experience with this product..."
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none"
            />
            <p className="text-sm text-gray-500 mt-1">
              {comment.length}/2000 characters
            </p>
          </div>

          {/* Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-6 py-3 border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting || orders.length === 0 || !selectedOrderId}
              className="flex-1 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-medium hover:from-purple-700 hover:to-pink-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? 'Submitting...' : 'Submit Review'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
