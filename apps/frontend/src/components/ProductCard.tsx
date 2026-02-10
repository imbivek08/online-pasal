import { Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { ShoppingCart, Star } from 'lucide-react';
import { api, type Product, type ProductRatingStats } from '../lib/api';
import { useCart } from '../contexts/CartContext';
import { useToast } from '../contexts/ToastContext';

interface ProductCardProps {
  product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
  const [adding, setAdding] = useState(false);
  const [ratingStats, setRatingStats] = useState<ProductRatingStats | null>(null);
  const { addToCart } = useCart();
  const toast = useToast();

  useEffect(() => {
    api.getProductRatingStats(product.id)
      .then(stats => setRatingStats(stats))
      .catch(() => {});
  }, [product.id]);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-NP', {
      style: 'currency',
      currency: 'NPR',
      minimumFractionDigits: 2,
    }).format(price);
  };

  return (
    <Link
      to={`/products/${product.id}`}
      className="group bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden hover:-translate-y-1"
    >
      {/* Product Image */}
      <div className="relative aspect-square bg-gradient-to-br from-primary/5 to-secondary/5 overflow-hidden">
        {product.image_url ? (
          <img
            src={product.image_url}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
            onError={(e) => {
              // Fallback if image fails to load
              e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="400" height="400"%3E%3Crect width="400" height="400" fill="%23f3f4f6"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-family="sans-serif" font-size="48" fill="%239ca3af"%3ENo Image%3C/text%3E%3C/svg%3E';
            }}
          />
        ) : (
          <div className="flex items-center justify-center h-full">
            <svg
              className="w-24 h-24 text-gray-300"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        )}
        
        {/* Stock Badge */}
        {product.stock_quantity <= 5 && product.stock_quantity > 0 && (
          <div className="absolute top-2 right-2 bg-yellow-500 text-white px-2 py-1 rounded-full text-xs font-semibold">
            Only {product.stock_quantity} left
          </div>
        )}
        {product.stock_quantity === 0 && (
          <div className="absolute top-2 right-2 bg-red-500 text-white px-2 py-1 rounded-full text-xs font-semibold">
            Out of Stock
          </div>
        )}
      </div>

      {/* Product Info */}
      <div className="p-4">
        <h3 className="font-semibold text-base text-gray-900 mb-1 line-clamp-2 group-hover:text-primary transition-colors">
          {product.name}
        </h3>
        
        {product.description && (
          <p className="text-xs text-gray-500 mb-3 line-clamp-2">
            {product.description}
          </p>
        )}

        {/* Rating */}
        <div className="flex items-center gap-1 mb-3">
          {ratingStats && ratingStats.total_reviews > 0 ? (
            <>
              <div className="flex items-center">
                {[1, 2, 3, 4, 5].map((star) => (
                  <Star
                    key={star}
                    className={`w-3.5 h-3.5 ${
                      star <= Math.round(ratingStats.average_rating)
                        ? 'fill-yellow-400 text-yellow-400'
                        : 'fill-gray-200 text-gray-200'
                    }`}
                  />
                ))}
              </div>
              <span className="text-xs text-gray-600 font-medium">
                {ratingStats.average_rating.toFixed(1)}
              </span>
              <span className="text-xs text-gray-400">
                ({ratingStats.total_reviews})
              </span>
            </>
          ) : (
            <span className="text-xs text-gray-400">No reviews yet</span>
          )}
        </div>

        <div className="flex items-center justify-between">
          <div>
            <div className="text-lg font-bold text-primary">
              {formatPrice(product.price)}
            </div>
          </div>
          
          <button
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-all flex items-center gap-1 ${
              product.stock_quantity > 0
                ? 'bg-primary text-white hover:bg-primary/90 hover:shadow-md'
                : 'bg-gray-200 text-gray-500 cursor-not-allowed'
            }`}
            disabled={product.stock_quantity === 0 || adding}
            onClick={async (e) => {
              e.preventDefault(); // Prevent navigation when clicking the button
              if (product.stock_quantity > 0) {
                setAdding(true);
                try {
                  await addToCart(product.id, 1);
                  // Optional: Show success toast/notification
                } catch (error) {
                  console.error('Failed to add to cart:', error);
                  toast.error('Failed to add to cart. Please try again.');
                } finally {
                  setAdding(false);
                }
              }
            }}
          >
            {adding ? (
              'Adding...'
            ) : (
              <>
                <ShoppingCart className="w-4 h-4" />
                {product.stock_quantity > 0 ? 'Add' : 'Out'}
              </>
            )}
          </button>
        </div>
      </div>
    </Link>
  );
}
