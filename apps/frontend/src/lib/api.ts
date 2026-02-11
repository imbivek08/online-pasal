import { useAuth } from '@clerk/clerk-react';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface User {
  id: string;
  clerk_id: string;
  email: string;
  username?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  avatar_url?: string;
  is_active: boolean;
  role: string;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

export interface Product {
  id: string;
  shop_id: string;
  category_id?: string;
  name: string;
  description?: string;
  price: number;
  stock_quantity: number;
  image_url?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProductSearchParams {
  search?: string;
  min_price?: number;
  max_price?: number;
  sort_by?: 'price_asc' | 'price_desc' | 'name_asc' | 'newest';
}

export interface BecomeVendorRequest {
  business_name: string;
  phone: string;
  business_description?: string;
}

export interface BecomeVendorResponse {
  id: string;
  email: string;
  role: string;
  phone: string;
  message: string;
  can_create_shop: boolean;
  next_step: string;
}

export interface Shop {
  id: string;
  vendor_id: string;
  name: string;
  slug: string;
  description: string;
  logo_url?: string;
  banner_url?: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
  phone?: string;
  email?: string;
  is_active: boolean;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateShopRequest {
  name: string;
  description: string;
  logo_url?: string;
  banner_url?: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
  phone?: string;
  email?: string;
}

export interface CartItem {
  id: string;
  product_id: string;
  product_name: string;
  product_price: number;
  product_image_url?: string;
  stock_quantity: number;
  is_active: boolean;
  shop_id: string;
  shop_name: string;
  quantity: number;
  subtotal: number;
  created_at: string;
  updated_at: string;
}

export interface Cart {
  id: string;
  user_id?: string;
  items: CartItem[];
  item_count: number;
  subtotal: number;
  created_at: string;
  updated_at: string;
}

export interface AddToCartRequest {
  product_id: string;
  quantity: number;
}

export interface UpdateCartItemRequest {
  quantity: number;
}

export interface Address {
  id: string;
  user_id: string;
  full_name: string;
  phone: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state?: string;
  postal_code?: string;
  country: string;
  is_default: boolean;
  address_type: string;
  created_at: string;
  updated_at: string;
}

export interface AddressInput {
  full_name: string;
  phone: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state?: string;
  postal_code?: string;
  country: string;
  is_default?: boolean;
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name: string;
  product_image_url?: string;
  shop_id: string;
  shop_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  created_at: string;
}

export interface Order {
  id: string;
  user_id: string;
  order_number: string;
  status: string;
  shipping_address?: Address;
  billing_address?: Address;
  items: OrderItem[];
  subtotal: number;
  shipping_cost: number;
  tax: number;
  discount: number;
  total: number;
  payment_method?: string;
  payment_status: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  confirmed_at?: string;
  shipped_at?: string;
  delivered_at?: string;
}

export interface CreateOrderRequest {
  shipping_address_id?: string;
  shipping_address?: AddressInput;
  billing_address?: AddressInput;
  payment_method: string;
  use_same_address: boolean;
  notes?: string;
}

export interface StripeCheckoutResponse {
  order: Order;
  checkout_url: string;
}

export interface StripeSessionStatus {
  session_id: string;
  payment_status: string;
  order_id: string;
  order_number: string;
}

export interface UpdateOrderStatusRequest {
  status: string;
}

export interface Review {
  id: string;
  product_id: string;
  user_id: string;
  order_id?: string;
  rating: number;
  title?: string;
  comment?: string;
  is_verified_purchase: boolean;
  is_approved: boolean;
  helpful_count: number;
  created_at: string;
  updated_at: string;
  user_name?: string;
  user_avatar?: string;
}

export interface CreateReviewRequest {
  product_id: string;
  order_id: string;
  rating: number;
  title?: string;
  comment?: string;
}

export interface UpdateReviewRequest {
  rating?: number;
  title?: string;
  comment?: string;
}

export interface ProductRatingStats {
  product_id: string;
  average_rating: number;
  total_reviews: number;
  five_star_count: number;
  four_star_count: number;
  three_star_count: number;
  two_star_count: number;
  one_star_count: number;
}

export interface ReviewListResponse {
  reviews: Review[];
  total_reviews: number;
  page: number;
  limit: number;
}

export interface CanReviewResponse {
  can_review: boolean;
  reason?: string;
  existing_review_id?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

class ApiClient {
  private getToken: (() => Promise<string | null>) | null = null;

  setAuth(getToken: () => Promise<string | null>) {
    this.getToken = getToken;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const token = this.getToken ? await this.getToken() : null;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || data.message || 'Request failed');
      }

      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Health check
  async healthCheck(): Promise<ApiResponse<{ status: string }>> {
    return this.request('/health');
  }

  // User endpoints
  async getProfile(): Promise<ApiResponse<User>> {
    return this.request('/api/v1/users/profile');
  }

  async updateProfile(data: Partial<User>): Promise<ApiResponse<User>> {
    return this.request('/api/v1/users/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteAccount(): Promise<ApiResponse<null>> {
    return this.request('/api/v1/users/account', {
      method: 'DELETE',
    });
  }

  async getUserById(id: string): Promise<ApiResponse<User>> {
    return this.request(`/api/v1/users/${id}`);
  }

  async becomeVendor(data: BecomeVendorRequest): Promise<ApiResponse<BecomeVendorResponse>> {
    return this.request('/api/v1/users/become-vendor', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getMyRole(): Promise<ApiResponse<{
    role: string;
    can_sell: boolean;
    can_buy: boolean;
    is_admin: boolean;
  }>> {
    return this.request('/api/v1/users/my-role');
  }

  // Product endpoints
  async getProducts(params?: ProductSearchParams): Promise<ApiResponse<Product[]>> {
    const queryParts: string[] = [];
    if (params?.search) queryParts.push(`search=${encodeURIComponent(params.search)}`);
    if (params?.min_price !== undefined) queryParts.push(`min_price=${params.min_price}`);
    if (params?.max_price !== undefined) queryParts.push(`max_price=${params.max_price}`);
    if (params?.sort_by) queryParts.push(`sort_by=${params.sort_by}`);
    const queryString = queryParts.length > 0 ? `?${queryParts.join('&')}` : '';
    return this.request(`/api/v1/products${queryString}`);
  }

  async getProductById(id: string): Promise<ApiResponse<Product>> {
    return this.request(`/api/v1/products/${id}`);
  }

  // Shop endpoints
  async getMyShop(): Promise<ApiResponse<Shop | null>> {
    return this.request('/api/v1/shops/my');
  }

  async createShop(data: CreateShopRequest): Promise<ApiResponse<Shop>> {
    return this.request('/api/v1/shops', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getShopById(id: string): Promise<ApiResponse<Shop>> {
    return this.request(`/api/v1/shops/${id}`);
  }

  async getShopBySlug(slug: string): Promise<ApiResponse<Shop>> {
    return this.request(`/api/v1/shops/slug/${slug}`);
  }

  // Cart endpoints
  async getCart(): Promise<ApiResponse<Cart>> {
    return this.request('/api/v1/cart');
  }

  async getCartItemCount(): Promise<ApiResponse<{ count: number }>> {
    return this.request('/api/v1/cart/count');
  }

  async addToCart(data: AddToCartRequest): Promise<ApiResponse<CartItem>> {
    return this.request('/api/v1/cart/items', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateCartItem(itemId: string, data: UpdateCartItemRequest): Promise<ApiResponse<null>> {
    return this.request(`/api/v1/cart/items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async removeCartItem(itemId: string): Promise<ApiResponse<null>> {
    return this.request(`/api/v1/cart/items/${itemId}`, {
      method: 'DELETE',
    });
  }

  async clearCart(): Promise<ApiResponse<null>> {
    return this.request('/api/v1/cart', {
      method: 'DELETE',
    });
  }

  // Order endpoints
  async createOrder(data: CreateOrderRequest): Promise<ApiResponse<Order>> {
    return this.request('/api/v1/orders', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getOrders(): Promise<ApiResponse<Order[]>> {
    return this.request('/api/v1/orders');
  }

  async getOrderById(orderId: string): Promise<ApiResponse<Order>> {
    return this.request(`/api/v1/orders/${orderId}`);
  }

  async cancelOrder(orderId: string): Promise<ApiResponse<null>> {
    return this.request(`/api/v1/orders/${orderId}/cancel`, {
      method: 'POST',
    });
  }

  // Stripe checkout
  async createStripeCheckout(data: CreateOrderRequest): Promise<ApiResponse<StripeCheckoutResponse>> {
    return this.request('/api/v1/orders/checkout/stripe', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async verifyStripeSession(sessionId: string): Promise<ApiResponse<StripeSessionStatus>> {
    return this.request(`/api/v1/orders/checkout/verify?session_id=${encodeURIComponent(sessionId)}`);
  }

  async getVendorOrders(): Promise<ApiResponse<Order[]>> {
    return this.request('/api/v1/vendor/orders');
  }

  async updateOrderStatus(orderId: string, data: UpdateOrderStatusRequest): Promise<ApiResponse<null>> {
    return this.request(`/api/v1/vendor/orders/${orderId}/status`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    });
  }

  // Address endpoints
  async getAddresses(): Promise<ApiResponse<Address[]>> {
    return this.request('/api/v1/addresses');
  }

  async getDefaultAddress(): Promise<ApiResponse<Address | null>> {
    return this.request('/api/v1/addresses/default');
  }

  async getAddressById(id: string): Promise<ApiResponse<Address>> {
    return this.request(`/api/v1/addresses/${id}`);
  }

  async createAddress(data: AddressInput): Promise<ApiResponse<Address>> {
    return this.request('/api/v1/addresses', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateAddress(id: string, data: AddressInput): Promise<ApiResponse<Address>> {
    return this.request(`/api/v1/addresses/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteAddress(id: string): Promise<ApiResponse<null>> {
    return this.request(`/api/v1/addresses/${id}`, {
      method: 'DELETE',
    });
  }

  async setDefaultAddress(id: string): Promise<ApiResponse<null>> {
    return this.request(`/api/v1/addresses/${id}/default`, {
      method: 'PATCH',
    });
  }

  // Review APIs
  async createReview(data: CreateReviewRequest): Promise<Review> {
    const response = await this.request('/api/v1/reviews', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data as Review;
  }

  async getProductReviews(productId: string, page: number = 1, limit: number = 10): Promise<ReviewListResponse> {
    const response = await this.request(`/api/v1/reviews/product/${productId}?page=${page}&limit=${limit}`);
    return response.data as ReviewListResponse;
  }

  async getProductRatingStats(productId: string): Promise<ProductRatingStats> {
    const response = await this.request(`/api/v1/reviews/product/${productId}/stats`);
    return response.data as ProductRatingStats;
  }

  async getReviewById(reviewId: string): Promise<Review> {
    const response = await this.request(`/api/v1/reviews/${reviewId}`);
    return response.data as Review;
  }

  async updateReview(reviewId: string, data: UpdateReviewRequest): Promise<void> {
    await this.request(`/api/v1/reviews/${reviewId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteReview(reviewId: string): Promise<void> {
    await this.request(`/api/v1/reviews/${reviewId}`, {
      method: 'DELETE',
    });
  }

  async markReviewHelpful(reviewId: string): Promise<void> {
    await this.request(`/api/v1/reviews/${reviewId}/helpful`, {
      method: 'POST',
    });
  }

  async canUserReviewProduct(productId: string): Promise<CanReviewResponse> {
    const response = await this.request(`/api/v1/reviews/can-review/${productId}`);
    return response.data as CanReviewResponse;
  }
}

export const api = new ApiClient();

// Custom hook to use API with auth
export function useApi() {
  const { getToken } = useAuth();
  
  // Set the token getter
  api.setAuth(getToken);
  
  return api;
}
