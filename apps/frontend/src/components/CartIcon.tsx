import { ShoppingCart } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useCart } from '../contexts/CartContext';

export default function CartIcon() {
  const { cartCount } = useCart();

  return (
    <Link
      to="/cart"
      className="relative inline-flex items-center justify-center p-2 rounded-lg hover:bg-gray-100 transition-colors"
    >
      <ShoppingCart className="w-6 h-6 text-gray-700" />
      {cartCount > 0 && (
        <span className="absolute -top-1 -right-1 bg-red-600 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
          {cartCount > 99 ? '99+' : cartCount}
        </span>
      )}
    </Link>
  );
}
