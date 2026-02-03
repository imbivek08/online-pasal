import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ShoppingBag, Trash2 } from 'lucide-react';
import { useCart } from '../contexts/CartContext';
import CartItemCard from '../components/CartItemCard';

export default function CartPage() {
  const { cart, loading, refreshCart, updateCartItem, removeCartItem, clearCart } = useCart();
  const navigate = useNavigate();

  useEffect(() => {
    refreshCart();
  }, []);

  const handleClearCart = async () => {
    if (!confirm('Are you sure you want to clear your entire cart?')) {
      return;
    }
    try {
      await clearCart();
    } catch (error) {
      console.error('Failed to clear cart:', error);
      alert('Failed to clear cart');
    }
  };

  const handleCheckout = () => {
    navigate('/checkout');
  };

  if (loading && !cart) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading your cart...</p>
        </div>
      </div>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center max-w-md px-4">
          <ShoppingBag className="w-24 h-24 text-gray-300 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Your cart is empty</h2>
          <p className="text-gray-600 mb-6">
            Looks like you haven't added any items to your cart yet.
          </p>
          <button
            onClick={() => navigate('/products')}
            className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 transition-colors"
          >
            Continue Shopping
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Shopping Cart</h1>
          <button
            onClick={handleClearCart}
            className="flex items-center gap-2 text-red-600 hover:text-red-800 transition-colors"
          >
            <Trash2 className="w-5 h-5" />
            Clear Cart
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items */}
          <div className="lg:col-span-2 space-y-4">
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <h2 className="text-lg font-semibold mb-4">
                Cart Items ({cart.item_count} {cart.item_count === 1 ? 'item' : 'items'})
              </h2>
            </div>
            
            {cart.items.map((item) => (
              <CartItemCard
                key={item.id}
                item={item}
                onUpdateQuantity={updateCartItem}
                onRemove={removeCartItem}
              />
            ))}
          </div>

          {/* Order Summary */}
          <div className="lg:col-span-1">
            <div className="bg-white p-6 rounded-lg border border-gray-200 sticky top-8">
              <h2 className="text-xl font-bold mb-4">Order Summary</h2>
              
              <div className="space-y-3 mb-6">
                <div className="flex justify-between text-gray-600">
                  <span>Subtotal</span>
                  <span>NPR {cart.subtotal.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Shipping</span>
                  <span className="text-green-600">Calculated at checkout</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Tax</span>
                  <span>NPR 0.00</span>
                </div>
                <div className="border-t pt-3">
                  <div className="flex justify-between text-xl font-bold">
                    <span>Total</span>
                    <span className="text-green-600">NPR {cart.subtotal.toFixed(2)}</span>
                  </div>
                </div>
              </div>

              <button
                onClick={handleCheckout}
                className="w-full bg-green-600 text-white py-3 rounded-lg hover:bg-green-700 transition-colors font-semibold"
              >
                Proceed to Checkout
              </button>

              <button
                onClick={() => navigate('/products')}
                className="w-full mt-3 border border-gray-300 text-gray-700 py-3 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Continue Shopping
              </button>

              <div className="mt-6 pt-6 border-t">
                <div className="flex items-start gap-2 text-sm text-gray-600">
                  <ShoppingBag className="w-5 h-5 flex-shrink-0 mt-0.5" />
                  <p>
                    We accept all major payment methods. Your transaction is secure and encrypted.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
