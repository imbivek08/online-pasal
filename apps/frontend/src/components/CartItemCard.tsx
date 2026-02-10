import { Minus, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import type { CartItem } from '../lib/api';
import { useToast } from '../contexts/ToastContext';

interface CartItemCardProps {
  item: CartItem;
  onUpdateQuantity: (itemId: string, quantity: number) => Promise<void>;
  onRemove: (itemId: string) => Promise<void>;
}

export default function CartItemCard({ item, onUpdateQuantity, onRemove }: CartItemCardProps) {
  const [updating, setUpdating] = useState(false);
  const [removing, setRemoving] = useState(false);
  const toast = useToast();

  const handleIncrease = async () => {
    if (item.quantity >= item.stock_quantity) {
      toast.warning('Cannot add more than available stock');
      return;
    }
    setUpdating(true);
    try {
      await onUpdateQuantity(item.id, item.quantity + 1);
    } catch (error) {
      console.error('Failed to update quantity:', error);
      toast.error('Failed to update quantity');
    } finally {
      setUpdating(false);
    }
  };

  const handleDecrease = async () => {
    if (item.quantity <= 1) {
      return;
    }
    setUpdating(true);
    try {
      await onUpdateQuantity(item.id, item.quantity - 1);
    } catch (error) {
      console.error('Failed to update quantity:', error);
      toast.error('Failed to update quantity');
    } finally {
      setUpdating(false);
    }
  };

  const handleRemove = async () => {
    setRemoving(true);
    try {
      await onRemove(item.id);
    } catch (error) {
      console.error('Failed to remove item:', error);
      toast.error('Failed to remove item');
    } finally {
      setRemoving(false);
    }
  };

  return (
    <div className="flex gap-4 p-4 bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition-shadow">
      {/* Product Image */}
      <div className="flex-shrink-0">
        <img
          src={item.product_image_url || 'https://via.placeholder.com/150'}
          alt={item.product_name}
          className="w-24 h-24 object-cover rounded-md"
        />
      </div>

      {/* Product Details */}
      <div className="flex-grow">
        <h3 className="font-semibold text-lg text-gray-900">{item.product_name}</h3>
        <p className="text-sm text-gray-600 mt-1">by {item.shop_name}</p>
        <p className="text-lg font-bold text-green-600 mt-2">
          NPR {item.product_price.toFixed(2)}
        </p>
        <p className="text-xs text-gray-500 mt-1">
          {item.stock_quantity > 0 ? (
            <span className="text-green-600">In Stock ({item.stock_quantity} available)</span>
          ) : (
            <span className="text-red-600">Out of Stock</span>
          )}
        </p>
      </div>

      {/* Quantity Controls */}
      <div className="flex flex-col items-end justify-between">
        <div className="flex items-center gap-2 border border-gray-300 rounded-lg">
          <button
            onClick={handleDecrease}
            disabled={updating || item.quantity <= 1}
            className="p-2 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Minus className="w-4 h-4" />
          </button>
          <span className="px-4 font-semibold">{item.quantity}</span>
          <button
            onClick={handleIncrease}
            disabled={updating || item.quantity >= item.stock_quantity}
            className="p-2 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Plus className="w-4 h-4" />
          </button>
        </div>

        {/* Subtotal and Remove */}
        <div className="text-right">
          <p className="text-xl font-bold text-gray-900">
            NPR {item.subtotal.toFixed(2)}
          </p>
          <button
            onClick={handleRemove}
            disabled={removing}
            className="mt-2 text-red-600 hover:text-red-800 flex items-center gap-1 text-sm disabled:opacity-50"
          >
            <Trash2 className="w-4 h-4" />
            Remove
          </button>
        </div>
      </div>
    </div>
  );
}
