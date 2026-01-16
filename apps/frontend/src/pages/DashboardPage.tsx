import { useUser } from '@clerk/clerk-react';
import { useEffect, useState } from 'react';
import { useApi, type User } from '../lib/api';
import Navbar from '../components/Navbar';

export default function DashboardPage() {
  const { user, isLoaded } = useUser();
  const api = useApi();
  const [profile, setProfile] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [healthStatus, setHealthStatus] = useState<string>('');

  useEffect(() => {
    if (isLoaded && user) {
      loadProfile();
      checkHealth();
    }
  }, [isLoaded, user]);

  const loadProfile = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.getProfile();
      setProfile(response.data || null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load profile');
    } finally {
      setLoading(false);
    }
  };

  const checkHealth = async () => {
    try {
      const response = await api.healthCheck();
      setHealthStatus(response.data?.status || 'unknown');
    } catch (err) {
      setHealthStatus('error');
    }
  };

  if (!isLoaded) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-xl text-gray-600">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <main className="pt-20 sm:pt-24 px-3 sm:px-4 md:px-6 lg:px-8 pb-8 sm:pb-12">
        <div className="max-w-7xl mx-auto">
          {/* Welcome Card */}
          <div className="bg-white rounded-xl sm:rounded-2xl shadow-lg p-4 sm:p-6 md:p-8 mb-6 sm:mb-8">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 sm:gap-0 mb-4 sm:mb-6">
              <h2 className="text-xl sm:text-2xl md:text-3xl font-bold text-gray-900">
                🎉 Welcome to Your Dashboard!
              </h2>
              <div className="flex items-center space-x-2 text-sm sm:text-base">
                <span className="text-gray-600">API Status:</span>
                <span className={`px-2 sm:px-3 py-1 rounded-full text-xs sm:text-sm font-medium ${
                  healthStatus === 'healthy' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-red-100 text-red-800'
                }`}>
                  {healthStatus || 'checking...'}
                </span>
              </div>
            </div>

            {loading ? (
              <div className="text-center py-6 sm:py-8 text-gray-600">Loading profile...</div>
            ) : error ? (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4 sm:p-6">
                <p className="text-sm sm:text-base text-red-800 mb-3 sm:mb-4">⚠️ {error}</p>
                <button 
                  onClick={loadProfile}
                  className="bg-primary text-white px-4 sm:px-6 py-2 rounded-lg font-medium hover:bg-primary-dark transition-colors text-sm sm:text-base w-full sm:w-auto"
                >
                  Retry
                </button>
                <p className="text-xs sm:text-sm text-red-600 mt-3 sm:mt-4">
                  Make sure your backend is running and you've set up the Clerk webhook.
                </p>
              </div>
            ) : profile ? (
              <div>
                <h3 className="text-lg sm:text-xl font-semibold text-gray-900 mb-3 sm:mb-4">Your Profile</h3>
                <div className="grid sm:grid-cols-2 gap-4 sm:gap-6">
                  <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                    <span className="text-xs sm:text-sm text-gray-600 block mb-1">Email:</span>
                    <span className="text-sm sm:text-base md:text-lg font-medium text-gray-900 break-all">{profile.email}</span>
                  </div>
                  <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                    <span className="text-xs sm:text-sm text-gray-600 block mb-1">Role:</span>
                    <span className="inline-block px-2 sm:px-3 py-1 bg-primary/10 text-primary rounded-full text-xs sm:text-sm font-medium">
                      {profile.role}
                    </span>
                  </div>
                  <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                    <span className="text-xs sm:text-sm text-gray-600 block mb-1">Status:</span>
                    <span className={`inline-block px-2 sm:px-3 py-1 rounded-full text-xs sm:text-sm font-medium ${
                      profile.is_active 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {profile.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </div>
                  <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                    <span className="text-xs sm:text-sm text-gray-600 block mb-1">Member Since:</span>
                    <span className="text-sm sm:text-base md:text-lg font-medium text-gray-900">
                      {new Date(profile.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  {profile.username && (
                    <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                      <span className="text-xs sm:text-sm text-gray-600 block mb-1">Username:</span>
                      <span className="text-sm sm:text-base md:text-lg font-medium text-gray-900">{profile.username}</span>
                    </div>
                  )}
                </div>
              </div>
            ) : null}
          </div>

          {/* Authentication Details */}
          <div className="bg-white rounded-xl sm:rounded-2xl shadow-lg p-4 sm:p-6 md:p-8 mb-6 sm:mb-8">
            <h3 className="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 mb-4 sm:mb-6">🔐 Authentication Details</h3>
            <div className="space-y-3 sm:space-y-4">
              <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                <span className="text-xs sm:text-sm text-gray-600 block mb-2">Clerk User ID:</span>
                <code className="text-xs sm:text-sm bg-gray-900 text-green-400 px-2 sm:px-3 py-1 rounded font-mono break-all block">
                  {user?.id}
                </code>
              </div>
              <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                <span className="text-xs sm:text-sm text-gray-600 block mb-2">Email:</span>
                <code className="text-xs sm:text-sm bg-gray-900 text-green-400 px-2 sm:px-3 py-1 rounded font-mono break-all block">
                  {user?.emailAddresses[0].emailAddress}
                </code>
              </div>
              <div className="bg-gray-50 p-3 sm:p-4 rounded-lg">
                <span className="text-xs sm:text-sm text-gray-600 block mb-2">Auth Provider:</span>
                <code className="text-xs sm:text-sm bg-gray-900 text-green-400 px-2 sm:px-3 py-1 rounded font-mono break-all block">
                  Clerk
                </code>
              </div>
            </div>
          </div>

          {/* Test API Endpoints */}
          <div className="bg-gradient-to-br from-primary/10 to-secondary/10 rounded-xl sm:rounded-2xl shadow-lg p-4 sm:p-6 md:p-8">
            <h3 className="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 mb-4 sm:mb-6">🧪 Test API Endpoints</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
              <button 
                onClick={loadProfile}
                className="bg-white hover:bg-gray-50 text-gray-900 px-4 sm:px-6 py-3 sm:py-4 rounded-xl font-medium shadow-md hover:shadow-lg transition-all text-sm sm:text-base"
              >
                🔄 Refresh Profile
              </button>
              <button 
                onClick={checkHealth}
                className="bg-white hover:bg-gray-50 text-gray-900 px-4 sm:px-6 py-3 sm:py-4 rounded-xl font-medium shadow-md hover:shadow-lg transition-all text-sm sm:text-base"
              >
                💚 Check API Health
              </button>
              <button 
                onClick={() => window.open('http://localhost:8080/health', '_blank')}
                className="bg-gradient-to-r from-primary to-secondary text-white px-4 sm:px-6 py-3 sm:py-4 rounded-xl font-medium shadow-md hover:shadow-lg transition-all text-sm sm:text-base sm:col-span-2 lg:col-span-1"
              >
                🌐 Open API Docs
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
