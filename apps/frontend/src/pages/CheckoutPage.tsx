import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { CreditCard, Truck, CheckCircle, MapPin, Home } from 'lucide-react';
import { useCart } from '../contexts/CartContext';
import { useToast } from '../contexts/ToastContext';
import { useApi } from '../lib/api';
import type { Address, AddressInput } from '../lib/api';

const emptyAddress: AddressInput = {
  full_name: '',
  phone: '',
  address_line1: '',
  address_line2: '',
  city: '',
  state: '',
  postal_code: '',
  country: 'Nepal',
  is_default: false,
};

export default function CheckoutPage() {
  const navigate = useNavigate();
  const { cart } = useCart();
  const toast = useToast();
  const api = useApi();

  const [loading, setLoading] = useState(false);
  const [addressLoading, setAddressLoading] = useState(true);

  // Address mode: 'saved' uses an existing address ID, 'new' sends inline form
  const [savedAddress, setSavedAddress] = useState<Address | null>(null);
  const [addressMode, setAddressMode] = useState<'saved' | 'new'>('new');
  const [formAddress, setFormAddress] = useState<AddressInput>(emptyAddress);

  // Payment & notes
  const [paymentMethod, setPaymentMethod] = useState('stripe');
  const [useSameAddress, setUseSameAddress] = useState(true);
  const [notes, setNotes] = useState('');

  // Load best available address on mount: default → latest → empty
  useEffect(() => {
    loadSavedAddress();
  }, []);

  const loadSavedAddress = async () => {
    setAddressLoading(true);
    try {
      // First try the default address
      const defaultResp = await api.getDefaultAddress();
      if (defaultResp.data) {
        setSavedAddress(defaultResp.data);
        setAddressMode('saved');
        populateFormFromAddress(defaultResp.data);
        setAddressLoading(false);
        return;
      }

      // No default — fall back to the latest address from the list
      const allResp = await api.getAddresses();
      const addresses = allResp.data || [];
      if (addresses.length > 0) {
        // Already sorted by is_default DESC, created_at DESC on the backend
        setSavedAddress(addresses[0]);
        setAddressMode('saved');
        populateFormFromAddress(addresses[0]);
        setAddressLoading(false);
        return;
      }

      // No addresses at all — new user, auto-default their first address
      setAddressMode('new');
      setFormAddress({ ...emptyAddress, is_default: true });
    } catch {
      setAddressMode('new');
      setFormAddress({ ...emptyAddress, is_default: true });
    } finally {
      setAddressLoading(false);
    }
  };

  const populateFormFromAddress = (addr: Address) => {
    setFormAddress({
      full_name: addr.full_name,
      phone: addr.phone,
      address_line1: addr.address_line1,
      address_line2: addr.address_line2 || '',
      city: addr.city,
      state: addr.state || '',
      postal_code: addr.postal_code || '',
      country: addr.country,
    });
  };

  const switchToSaved = () => {
    if (!savedAddress) return;
    setAddressMode('saved');
    populateFormFromAddress(savedAddress);
  };

  const switchToNew = () => {
    setAddressMode('new');
    // If no saved address exists, this is the user's first address — auto-default it
    setFormAddress({ ...emptyAddress, is_default: !savedAddress });
  };

  const handleSubmitOrder = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!cart || cart.items.length === 0) {
      toast.warning('Your cart is empty');
      return;
    }

    setLoading(true);
    try {
      const orderPayload: any = {
        payment_method: paymentMethod,
        use_same_address: useSameAddress,
        notes: notes || undefined,
      };

      if (addressMode === 'saved' && savedAddress) {
        orderPayload.shipping_address_id = savedAddress.id;
      } else {
        if (!formAddress.full_name || !formAddress.phone || !formAddress.address_line1 || !formAddress.city || !formAddress.country) {
          toast.warning('Please fill in all required address fields');
          setLoading(false);
          return;
        }
        orderPayload.shipping_address = formAddress;
      }

      if (paymentMethod === 'stripe') {
        // Stripe flow: create checkout session and redirect
        const response = await api.createStripeCheckout(orderPayload);
        if (response.success && response.data?.checkout_url) {
          // Redirect to Stripe Checkout
          window.location.href = response.data.checkout_url;
          return;
        }
      } else {
        // COD flow: create order directly
        const response = await api.createOrder(orderPayload);
        if (response.success && response.data) {
          toast.success(`Order placed successfully! Order #${response.data.order_number}`);
          navigate(`/orders/${response.data.id}`);
        }
      }
    } catch (error: any) {
      console.error('Failed to create order:', error);
      toast.error(error.message || 'Failed to place order. Please try again.');
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
            <div className="flex items-center text-green-600">
              <Truck className="w-6 h-6" />
              <span className="ml-2 font-medium">Shipping</span>
            </div>
            <div className="w-16 h-1 mx-4 bg-gray-300"></div>
            <div className="flex items-center text-gray-400">
              <CreditCard className="w-6 h-6" />
              <span className="ml-2 font-medium">Payment</span>
            </div>
            <div className="w-16 h-1 mx-4 bg-gray-300"></div>
            <div className="flex items-center text-gray-400">
              <CheckCircle className="w-6 h-6" />
              <span className="ml-2 font-medium">Review</span>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmitOrder}>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-6">
              {/* Shipping Address */}
              <div className="bg-white p-6 rounded-lg border border-gray-200">
                <h2 className="text-xl font-bold mb-4">Shipping Address</h2>

                {addressLoading ? (
                  <div className="text-center py-8 text-gray-500">Loading address...</div>
                ) : (
                  <>
                    {/* Mode Toggle Buttons */}
                    {savedAddress && (
                      <div className="flex gap-3 mb-5">
                        <button
                          type="button"
                          onClick={switchToSaved}
                          className={`flex-1 flex items-center justify-center gap-2.5 p-4 rounded-lg border-2 transition-all ${
                            addressMode === 'saved'
                              ? 'border-green-500 bg-green-50 text-green-700'
                              : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:bg-gray-50'
                          }`}
                        >
                          <Home className="w-5 h-5" />
                          <div className="text-left">
                            <span className="font-medium text-sm block">{savedAddress.is_default ? 'Default Address' : 'Last Shipped Address'}</span>
                            <span className="text-xs opacity-70">{savedAddress.address_line1}, {savedAddress.city}</span>
                          </div>
                        </button>
                        <button
                          type="button"
                          onClick={switchToNew}
                          className={`flex-1 flex items-center justify-center gap-2.5 p-4 rounded-lg border-2 transition-all ${
                            addressMode === 'new'
                              ? 'border-green-500 bg-green-50 text-green-700'
                              : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:bg-gray-50'
                          }`}
                        >
                          <MapPin className="w-5 h-5" />
                          <span className="font-medium text-sm">Ship to New Location</span>
                        </button>
                      </div>
                    )}

                    {/* Address Form (always visible, pre-filled or empty based on mode) */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Full Name *</label>
                        <input
                          type="text"
                          required
                          value={formAddress.full_name}
                          onChange={(e) => setFormAddress({ ...formAddress, full_name: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Phone *</label>
                        <input
                          type="tel"
                          required
                          value={formAddress.phone}
                          onChange={(e) => setFormAddress({ ...formAddress, phone: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-gray-700 mb-1">Address Line 1 *</label>
                        <input
                          type="text"
                          required
                          value={formAddress.address_line1}
                          onChange={(e) => setFormAddress({ ...formAddress, address_line1: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          placeholder="Street address, house number"
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-gray-700 mb-1">Address Line 2</label>
                        <input
                          type="text"
                          value={formAddress.address_line2}
                          onChange={(e) => setFormAddress({ ...formAddress, address_line2: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          placeholder="Apartment, suite, unit, etc."
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">City *</label>
                        <input
                          type="text"
                          required
                          value={formAddress.city}
                          onChange={(e) => setFormAddress({ ...formAddress, city: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">State / Province</label>
                        <input
                          type="text"
                          value={formAddress.state}
                          onChange={(e) => setFormAddress({ ...formAddress, state: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
                        <input
                          type="text"
                          value={formAddress.postal_code}
                          onChange={(e) => setFormAddress({ ...formAddress, postal_code: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Country *</label>
                        <input
                          type="text"
                          required
                          value={formAddress.country}
                          onChange={(e) => setFormAddress({ ...formAddress, country: e.target.value })}
                          readOnly={addressMode === 'saved'}
                          className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent ${addressMode === 'saved' ? 'bg-gray-50 text-gray-600 cursor-not-allowed' : ''}`}
                        />
                      </div>
                    </div>

                    {/* Save as default checkbox — only in 'new' mode */}
                    {addressMode === 'new' && (
                      <div className="mt-4">
                        <label className="flex items-center cursor-pointer">
                          <input
                            type="checkbox"
                            checked={formAddress.is_default ?? false}
                            onChange={(e) => setFormAddress({ ...formAddress, is_default: e.target.checked })}
                            className="w-4 h-4 text-green-600 rounded focus:ring-green-500"
                          />
                          <span className="ml-2 text-sm text-gray-700">
                            {savedAddress ? 'Save as my default address' : 'This is your first address — it will be saved as default'}
                          </span>
                        </label>
                      </div>
                    )}
                  </>
                )}

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
              <div className="bg-white p-6 rounded-lg border border-gray-200">
                <h2 className="text-xl font-bold mb-4">Payment Method</h2>
                <div className="space-y-3">
                  <label className="flex items-center p-4 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50">
                    <input
                      type="radio"
                      name="payment"
                      value="stripe"
                      checked={paymentMethod === 'stripe'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="w-4 h-4 text-green-600"
                    />
                    <div className="ml-3">
                      <span className="font-medium">Pay with Card</span>
                      <p className="text-sm text-gray-500">Secure payment via Stripe (Visa, Mastercard, etc.)</p>
                    </div>
                  </label>
                  <label className="flex items-center p-4 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50">
                    <input
                      type="radio"
                      name="payment"
                      value="cash_on_delivery"
                      checked={paymentMethod === 'cash_on_delivery'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="w-4 h-4 text-green-600"
                    />
                    <div className="ml-3">
                      <span className="font-medium">Cash on Delivery</span>
                      <p className="text-sm text-gray-500">Pay when your order is delivered</p>
                    </div>
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
              <div className="bg-white p-6 rounded-lg border border-gray-200 sticky top-20">
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
                  {loading
                    ? 'Processing...'
                    : paymentMethod === 'stripe'
                      ? 'Pay with Card'
                      : 'Place Order'}
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
