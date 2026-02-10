import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ShoppingCart, ArrowLeft, Package, CheckCircle, Star } from 'lucide-react';
import { useAuth } from '@clerk/clerk-react';
import { api } from '../lib/api';
import type { Product, ProductRatingStats, CreateReviewRequest, Review } from '../lib/api';
import { useCart } from '../contexts/CartContext';
import ProductCard from '../components/ProductCard';
import ProductRatingSummary from '../components/ProductRatingSummary';
import ReviewCard from '../components/ReviewCard';
import WriteReviewModal from '../components/WriteReviewModal';

export default function ProductDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { getToken, isSignedIn } = useAuth();
  const { addToCart } = useCart();
  
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [quantity, setQuantity] = useState(1);
  const [addedToCart, setAddedToCart] = useState(false);
  const [relatedProducts, setRelatedProducts] = useState<Product[]>([]);
  
  // Review states
  const [reviews, setReviews] = useState<Review[]>([]);
  const [ratingStats, setRatingStats] = useState<ProductRatingStats | null>(null);
  const [reviewsPage, setReviewsPage] = useState(1);
  const [totalReviews, setTotalReviews] = useState(0);
  const [loadingReviews, setLoadingReviews] = useState(false);
  const [showWriteReviewModal, setShowWriteReviewModal] = useState(false);
  const [selectedOrderId, setSelectedOrderId] = useState<string>('');

  useEffect(() => {
    const fetchProduct = async () => {
      if (!id) {
        setError('Product ID is required');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        const token = await getToken();
        if (token) {
          api.setAuth(getToken);
        }
        
        const response = await api.getProductById(id);
        if (response.data) {
          setProduct(response.data);
          
          // Fetch all products to get related ones
          const productsResponse = await api.getProducts();
          if (productsResponse.success && productsResponse.data) {
            // Filter out current product and get up to 4 related products
            const related = productsResponse.data
              .filter(p => p.id !== id && p.is_active)
              .slice(0, 4);
            setRelatedProducts(related);
          }
        } else {
          setError('Product not found');
        }
      } catch (err) {
        console.error('Failed to fetch product:', err);
        setError('Failed to load product details');
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [id, getToken]);

  // Fetch reviews and rating stats
  useEffect(() => {
    const fetchReviews = async () => {
      if (!id) return;

      try {
        setLoadingReviews(true);
        
        // Fetch reviews
        const reviewsData = await api.getProductReviews(id, reviewsPage, 10);
        setReviews(reviewsData.reviews);
        setTotalReviews(reviewsData.total_reviews);

        // Fetch rating stats
        const stats = await api.getProductRatingStats(id);
        setRatingStats(stats);
      } catch (error) {
        console.error('Failed to fetch reviews:', error);
      } finally {
        setLoadingReviews(false);
      }
    };

    fetchReviews();
  }, [id, reviewsPage]);

  const handleWriteReview = async (reviewData: CreateReviewRequest) => {
    try {
      const token = await getToken();
      if (token) {
        api.setAuth(getToken);
      }
      
      await api.createReview(reviewData);
      
      // Refresh reviews after submission
      if (id) {
        const reviewsData = await api.getProductReviews(id, 1, 10);
        setReviews(reviewsData.reviews);
        setTotalReviews(reviewsData.total_reviews);

        const stats = await api.getProductRatingStats(id);
        setRatingStats(stats);
      }
    } catch (error) {
      throw error;
    }
  };

  const handleMarkHelpful = async (reviewId: string) => {
    try {
      const token = await getToken();
      if (token) {
        api.setAuth(getToken);
        await api.markReviewHelpful(reviewId);
        
        // Refresh reviews to show updated helpful count
        if (id) {
          const reviewsData = await api.getProductReviews(id, reviewsPage, 10);
          setReviews(reviewsData.reviews);
        }
      }
    } catch (error) {
      console.error('Failed to mark review as helpful:', error);
    }
  };

  const handleOpenWriteReview = () => {
    // For testing, use a hardcoded order ID
    // In production, you should fetch user's orders and let them select
    setSelectedOrderId('');
    setShowWriteReviewModal(true);
  };

  const handleAddToCart = async () => {
    if (!product || product.stock_quantity === 0) return;

    setAdding(true);
    try {
      await addToCart(product.id, quantity);
      setAddedToCart(true);
      setTimeout(() => setAddedToCart(false), 2000);
    } catch (error) {
      console.error('Failed to add to cart:', error);
      alert('Failed to add to cart. Please try again.');
    } finally {
      setAdding(false);
    }
  };

  const handleQuantityChange = (newQuantity: number) => {
    if (product && newQuantity >= 1 && newQuantity <= product.stock_quantity) {
      setQuantity(newQuantity);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-NP', {
      style: 'currency',
      currency: 'NPR',
      minimumFractionDigits: 2,
    }).format(price);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-primary mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading product details...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="text-red-500 text-6xl mb-4">⚠️</div>
            <h2 className="text-2xl font-bold text-gray-800 mb-2">Product Not Found</h2>
            <p className="text-gray-600 mb-6">{error || 'The product you are looking for does not exist.'}</p>
            <button
              onClick={() => navigate('/products')}
              className="bg-primary text-white px-6 py-2 rounded-lg hover:bg-primary/90 transition-colors"
            >
              Back to Products
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Back Button */}
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-2 text-gray-600 hover:text-primary transition-colors mb-6 group"
        >
          <ArrowLeft className="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
          <span>Back</span>
        </button>

        {/* Product Detail */}
        <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 p-6 lg:p-10">
            {/* Product Image */}
            <div className="relative h-96 bg-gradient-to-br from-primary/5 to-secondary/5 rounded-xl overflow-hidden">
              {product.image_url ? (
                <img
                  src={product.image_url}
                  alt={product.name}
                  className="w-full h-full object-cover"
                  onError={(e) => {
                    e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="600" height="600"%3E%3Crect width="600" height="600" fill="%23f3f4f6"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-family="sans-serif" font-size="64" fill="%239ca3af"%3ENo Image%3C/text%3E%3C/svg%3E';
                  }}
                />
              ) : (
                <div className="flex items-center justify-center h-full">
                  <Package className="w-32 h-32 text-gray-300" />
                </div>
              )}
              
              {/* Stock Badge */}
              {product.stock_quantity <= 5 && product.stock_quantity > 0 && (
                <div className="absolute top-4 right-4 bg-yellow-500 text-white px-4 py-2 rounded-full text-sm font-semibold shadow-lg">
                  Only {product.stock_quantity} left
                </div>
              )}
              {product.stock_quantity === 0 && (
                <div className="absolute top-4 right-4 bg-red-500 text-white px-4 py-2 rounded-full text-sm font-semibold shadow-lg">
                  Out of Stock
                </div>
              )}
            </div>

            {/* Product Info */}
            <div className="flex flex-col justify-between">
              <div>
                <h1 className="text-2xl lg:text-3xl font-bold text-gray-900 mb-4">
                  {product.name}
                </h1>

                {/* Price */}
                <div className="mb-6">
                  <div className="text-2xl font-bold text-primary">
                    {formatPrice(product.price)}
                  </div>
                </div>

                {/* Stock Status */}
                <div className="flex items-center gap-2 mb-6">
                  {product.stock_quantity > 0 ? (
                    <>
                      <CheckCircle className="w-5 h-5 text-green-500" />
                      <span className="text-green-600 font-medium">
                        In Stock
                      </span>
                    </>
                  ) : (
                    <>
                      <div className="w-5 h-5 rounded-full bg-red-500" />
                      <span className="text-red-600 font-medium">Out of Stock</span>
                    </>
                  )}
                </div>

                {/* Description */}
                {product.description && (
                  <div className="mb-8">
                    <h2 className="text-lg font-semibold text-gray-900 mb-2">Description</h2>
                    <p className="text-gray-600 leading-relaxed whitespace-pre-wrap">
                      {product.description}
                    </p>
                  </div>
                )}

              </div>

              {/* Add to Cart Section */}
              <div className="border-t border-gray-200 pt-6">
                {product.stock_quantity > 0 && (
                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Quantity
                    </label>
                    <div className="flex items-center gap-3">
                      <button
                        onClick={() => handleQuantityChange(quantity - 1)}
                        disabled={quantity <= 1}
                        className="w-10 h-10 rounded-lg border border-gray-300 flex items-center justify-center hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                      >
                        -
                      </button>
                      <input
                        type="number"
                        min="1"
                        max={product.stock_quantity}
                        value={quantity}
                        onChange={(e) => handleQuantityChange(parseInt(e.target.value) || 1)}
                        className="w-20 h-10 text-center border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                      <button
                        onClick={() => handleQuantityChange(quantity + 1)}
                        disabled={quantity >= product.stock_quantity}
                        className="w-10 h-10 rounded-lg border border-gray-300 flex items-center justify-center hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                      >
                        +
                      </button>
                    </div>
                  </div>
                )}

                <div className="flex gap-3">
                  <button
                    onClick={handleAddToCart}
                    disabled={product.stock_quantity === 0 || adding}
                    className={`flex-1 px-6 py-3 rounded-lg font-medium text-base flex items-center justify-center gap-2 transition-all ${
                      product.stock_quantity > 0
                        ? addedToCart
                          ? 'bg-green-500 text-white'
                          : 'bg-primary text-white hover:bg-primary/90 hover:shadow-lg'
                        : 'bg-gray-200 text-gray-500 cursor-not-allowed'
                    }`}
                  >
                    {adding ? (
                      <>
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                        Adding...
                      </>
                    ) : addedToCart ? (
                      <>
                        <CheckCircle className="w-4 h-4" />
                        Added to Cart!
                      </>
                    ) : (
                      <>
                        <ShoppingCart className="w-4 h-4" />
                        Add to Cart
                      </>
                    )}
                  </button>

                  {product.stock_quantity > 0 && (
                    <Link
                      to="/cart"
                      className="px-6 py-3 rounded-lg font-medium text-base border-2 border-primary text-primary hover:bg-primary hover:text-white transition-all"
                    >
                      View Cart
                    </Link>
                  )}
                </div>

                {product.stock_quantity === 0 && (
                  <p className="text-center text-gray-500 mt-4">
                    This product is currently unavailable
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Reviews Section */}
        <div className="mt-12">
          <div className="bg-white rounded-2xl shadow-lg p-6 lg:p-10">
            {/* Rating Summary */}
            {ratingStats && <ProductRatingSummary stats={ratingStats} />}

            {/* Write Review Button */}
            {isSignedIn && (
              <div className="mb-8">
                <button
                  onClick={handleOpenWriteReview}
                  className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-medium hover:from-purple-700 hover:to-pink-700 transition-all"
                >
                  <Star className="w-5 h-5" />
                  Write a Review
                </button>
              </div>
            )}

            {/* Reviews List */}
            <div>
              <h3 className="text-xl font-semibold mb-6">Customer Reviews</h3>
              
              {loadingReviews ? (
                <div className="flex justify-center py-12">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                </div>
              ) : reviews.length > 0 ? (
                <>
                  {reviews.map((review) => (
                    <ReviewCard
                      key={review.id}
                      review={review}
                      onMarkHelpful={handleMarkHelpful}
                    />
                  ))}
                  
                  {/* Pagination */}
                  {totalReviews > 10 && (
                    <div className="flex justify-center gap-2 mt-8">
                      <button
                        onClick={() => setReviewsPage(Math.max(1, reviewsPage - 1))}
                        disabled={reviewsPage === 1}
                        className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Previous
                      </button>
                      <span className="px-4 py-2 text-gray-700">
                        Page {reviewsPage} of {Math.ceil(totalReviews / 10)}
                      </span>
                      <button
                        onClick={() => setReviewsPage(reviewsPage + 1)}
                        disabled={reviewsPage >= Math.ceil(totalReviews / 10)}
                        className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Next
                      </button>
                    </div>
                  )}
                </>
              ) : (
                <div className="text-center py-12">
                  <div className="text-gray-400 text-6xl mb-4">💬</div>
                  <p className="text-gray-600 mb-4">No reviews yet</p>
                  {isSignedIn && (
                    <button
                      onClick={handleOpenWriteReview}
                      className="text-primary hover:text-primary/80 font-medium"
                    >
                      Be the first to review this product
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Related Products */}
        {relatedProducts.length > 0 && (
          <div className="mt-12">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Related Products</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {relatedProducts.map((relatedProduct) => (
                <ProductCard key={relatedProduct.id} product={relatedProduct} />
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Write Review Modal */}
      {showWriteReviewModal && product && (
        <WriteReviewModal
          productId={product.id}
          productName={product.name}
          orderId={selectedOrderId}
          onClose={() => setShowWriteReviewModal(false)}
          onSubmit={handleWriteReview}
        />
      )}
    </div>
  );
}
