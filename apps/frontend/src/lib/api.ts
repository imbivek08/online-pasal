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
}

export const api = new ApiClient();

// Custom hook to use API with auth
export function useApi() {
  const { getToken } = useAuth();
  
  // Set the token getter
  api.setAuth(getToken);
  
  return api;
}
