import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { useAuth } from '@clerk/clerk-react';
import { useApi } from '../lib/api';
import type { Cart } from '../lib/api';

interface CartContextType {
  cart: Cart | null;
  cartCount: number;
  loading: boolean;
  refreshCart: () => Promise<void>;
  addToCart: (productId: string, quantity: number) => Promise<void>;
  updateCartItem: (itemId: string, quantity: number) => Promise<void>;
  removeCartItem: (itemId: string) => Promise<void>;
  clearCart: () => Promise<void>;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: ReactNode }) {
  const [cart, setCart] = useState<Cart | null>(null);
  const [cartCount, setCartCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const { isSignedIn } = useAuth();
  const api = useApi();

  const refreshCart = async () => {
    if (!isSignedIn) {
      setCart(null);
      setCartCount(0);
      return;
    }

    try {
      setLoading(true);
      const response = await api.getCart();
      if (response.success && response.data) {
        setCart(response.data);
        setCartCount(response.data.item_count);
      }
    } catch (error) {
      console.error('Failed to fetch cart:', error);
      setCart(null);
      setCartCount(0);
    } finally {
      setLoading(false);
    }
  };

  const addToCart = async (productId: string, quantity: number) => {
    try {
      setLoading(true);
      const response = await api.addToCart({ product_id: productId, quantity });
      if (response.success) {
        await refreshCart();
      }
    } catch (error) {
      console.error('Failed to add to cart:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const updateCartItem = async (itemId: string, quantity: number) => {
    try {
      setLoading(true);
      const response = await api.updateCartItem(itemId, { quantity });
      if (response.success) {
        await refreshCart();
      }
    } catch (error) {
      console.error('Failed to update cart item:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const removeCartItem = async (itemId: string) => {
    try {
      setLoading(true);
      const response = await api.removeCartItem(itemId);
      if (response.success) {
        await refreshCart();
      }
    } catch (error) {
      console.error('Failed to remove cart item:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const clearCart = async () => {
    try {
      setLoading(true);
      const response = await api.clearCart();
      if (response.success) {
        await refreshCart();
      }
    } catch (error) {
      console.error('Failed to clear cart:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Fetch cart when user signs in
  useEffect(() => {
    if (isSignedIn) {
      refreshCart();
    } else {
      setCart(null);
      setCartCount(0);
    }
  }, [isSignedIn]);

  return (
    <CartContext.Provider
      value={{
        cart,
        cartCount,
        loading,
        refreshCart,
        addToCart,
        updateCartItem,
        removeCartItem,
        clearCart,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
}
