import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { CreditCard, Truck, CheckCircle } from 'lucide-react';
import { useCart } from '../contexts/CartContext';
import { useApi } from '../lib/api';
import type { AddressInput } from '../lib/api';

export default function CheckoutPage() {
  const navigate = useNavigate();
  const { cart } = useCart();
  const api = useApi();

  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  
  // Shipping address form
  const [shippingAddress, setShippingAddress] = useState<AddressInput>({
    full_name: '',
    phone: '',
    address_line1: '',
    address_line2: '',
    city: '',
    state: '',
    postal_code: '',
    country: 'Nepal',
  });

  // Payment details
  const [paymentMethod, setPaymentMethod] = useState('cash_on_delivery');
  const [useSameAddress, setUseSameAddress] = useState(true);
  const [notes, setNotes] = useState('');

  const handleSubmitOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!cart || cart.items.length === 0) {
      alert('Your cart is empty');
      return;
    }

    setLoading(true);
    try {
      const response = await api.createOrder({
        shipping_address: shippingAddress,
        payment_method: paymentMethod,
        use_same_address: useSameAddress,
        notes: notes || undefined,
      });

      if (response.success && response.data) {
        alert('Order placed successfully!');
        navigate(`/orders/${response.data.id}`);
      }
    } catch (error: any) {
      console.error('Failed to create order:', error);
      alert(error.message || 'Failed to place order. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Your cart is empty</h2>
          <p className="text-gray-600 mb-6">Add some items to your cart before checkout</p>
          <button
            onClick={() => navigate('/products')}
            className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700"
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
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Checkout</h1>

        {/* Progress Steps */}
        <div className="mb-8 flex items-center justify-center">
          <div className="flex items-center">
            <div className={`flex items-center ${step >= 1 ? 'text-green-600' : 'text-gray-400'}`}>
              <Truck className="w-6 h-6" />
              <span className="ml-2 font-medium">Shipping</span>
            </div>
            <div className="w-16 h-1 mx-4 bg-gray-300"></div>
            <div className={`flex items-center ${step >= 2 ? 'text-green-600' : 'text-gray-400'}`}>
              <CreditCard className="w-6 h-6" />
              <span className="ml-2 font-medium">Payment</span>
            </div>
            <div className="w-16 h-1 mx-4 bg-gray-300"></div>
            <div className={`flex items-center ${step >= 3 ? 'text-green-600' : 'text-gray-400'}`}>
              <CheckCircle className="w-6 h-6" />
              <span className="ml-2 font-medium">Review</span>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmitOrder}>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Main Content */}
            <div className="lg:col-span-2">
              {/* Shipping Address */}
              <div className="bg-white p-6 rounded-lg border border-gray-200 mb-6">
                <h2 className="text-xl font-bold mb-4">Shipping Address</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Full Name *</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.full_name}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, full_name: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Phone *</label>
                    <input
                      type="tel"
                      required
                      value={shippingAddress.phone}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, phone: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 mb-1">Address Line 1 *</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.address_line1}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, address_line1: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700 mb-1">Address Line 2</label>
                    <input
                      type="text"
                      value={shippingAddress.address_line2}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, address_line2: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">City *</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.city}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, city: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">State/Province</label>
                    <input
                      type="text"
                      value={shippingAddress.state}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, state: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
                    <input
                      type="text"
                      value={shippingAddress.postal_code}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, postal_code: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Country *</label>
                    <input
                      type="text"
                      required
                      value={shippingAddress.country}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, country: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                    />
                  </div>
                </div>

                <div className="mt-4">
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      checked={useSameAddress}
                      onChange={(e) => setUseSameAddress(e.target.checked)}
                      className="w-4 h-4 text-green-600 rounded focus:ring-green-500"
                    />
                    <span className="ml-2 text-sm text-gray-700">Billing address same as shipping</span>
                  </label>
                </div>
              </div>

              {/* Payment Method */}
              <div className="bg-white p-6 rounded-lg border border-gray-200 mb-6">
                <h2 className="text-xl font-bold mb-4">Payment Method</h2>
                <div className="space-y-3">
                  <label className="flex items-center p-4 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50">
                    <input
                      type="radio"
                      name="payment"
                      value="cash_on_delivery"
                      checked={paymentMethod === 'cash_on_delivery'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="w-4 h-4 text-green-600"
                    />
                    <span className="ml-3 font-medium">Cash on Delivery</span>
                  </label>
                  <label className="flex items-center p-4 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50">
                    <input
                      type="radio"
                      name="payment"
                      value="esewa"
                      checked={paymentMethod === 'esewa'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="w-4 h-4 text-green-600"
                    />
                    <span className="ml-3 font-medium">eSewa</span>
                  </label>
                  <label className="flex items-center p-4 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50">
                    <input
                      type="radio"
                      name="payment"
                      value="khalti"
                      checked={paymentMethod === 'khalti'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="w-4 h-4 text-green-600"
                    />
                    <span className="ml-3 font-medium">Khalti</span>
                  </label>
                </div>
              </div>

              {/* Order Notes */}
              <div className="bg-white p-6 rounded-lg border border-gray-200">
                <h2 className="text-xl font-bold mb-4">Order Notes (Optional)</h2>
                <textarea
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="Any special instructions for delivery?"
                  rows={4}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>

            {/* Order Summary */}
            <div className="lg:col-span-1">
              <div className="bg-white p-6 rounded-lg border border-gray-200 sticky top-8">
                <h2 className="text-xl font-bold mb-4">Order Summary</h2>
                
                <div className="space-y-3 mb-6 max-h-60 overflow-y-auto">
                  {cart.items.map((item) => (
                    <div key={item.id} className="flex gap-3">
                      <img
                        src={item.product_image_url || 'https://via.placeholder.com/80'}
                        alt={item.product_name}
                        className="w-16 h-16 object-cover rounded"
                      />
                      <div className="flex-1">
                        <p className="font-medium text-sm">{item.product_name}</p>
                        <p className="text-sm text-gray-600">Qty: {item.quantity}</p>
                        <p className="text-sm font-semibold text-green-600">NPR {item.subtotal.toFixed(2)}</p>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="space-y-2 mb-6 pt-4 border-t">
                  <div className="flex justify-between text-gray-600">
                    <span>Subtotal</span>
                    <span>NPR {cart.subtotal.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Shipping</span>
                    <span className="text-green-600">Calculated at delivery</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Tax</span>
                    <span>NPR 0.00</span>
                  </div>
                  <div className="flex justify-between text-xl font-bold pt-2 border-t">
                    <span>Total</span>
                    <span className="text-green-600">NPR {cart.subtotal.toFixed(2)}</span>
                  </div>
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-green-600 text-white py-3 rounded-lg hover:bg-green-700 transition-colors font-semibold disabled:opacity-50"
                >
                  {loading ? 'Placing Order...' : 'Place Order'}
                </button>

                <button
                  type="button"
                  onClick={() => navigate('/cart')}
                  className="w-full mt-3 border border-gray-300 text-gray-700 py-3 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  Back to Cart
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
