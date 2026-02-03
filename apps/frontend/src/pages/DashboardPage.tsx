import { useUser, useAuth } from '@clerk/clerk-react';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useApi, type User, type Shop } from '../lib/api';
import BecomeVendorModal from '../components/BecomeVendorModal';

export default function DashboardPage() {
  const { user, isLoaded } = useUser();
  const { getToken } = useAuth();
  const navigate = useNavigate();
  const api = useApi();
  const [profile, setProfile] = useState<User | null>(null);
  const [shop, setShop] = useState<Shop | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [healthStatus, setHealthStatus] = useState<string>('');
  const [jwtToken, setJwtToken] = useState<string>('');
  const [tokenCopied, setTokenCopied] = useState(false);
  const [showVendorModal, setShowVendorModal] = useState(false);

  useEffect(() => {
    if (isLoaded && user) {
      loadProfile();
      checkHealth();
      loadJwtToken();
    }
  }, [isLoaded, user]);

  useEffect(() => {
    // Load shop when profile is loaded and user is a vendor
    if (profile && profile.role === 'vendor') {
      loadShop();
    }
  }, [profile]);

  const loadJwtToken = async () => {
    try {
      const token = await getToken();
      if (token) {
        setJwtToken(token);
      }
    } catch (err) {
      console.error('Failed to get JWT token:', err);
    }
  };

  const copyTokenToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(jwtToken);
      setTokenCopied(true);
      setTimeout(() => setTokenCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy token:', err);
    }
  };

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

  const loadShop = async () => {
    try {
      const response = await api.getMyShop();
      setShop(response.data || null);
    } catch (err) {
      console.error('Failed to load shop:', err);
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

  const handleVendorSuccess = () => {
    // Reload profile to get updated role
    loadProfile();
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

          {/* Become Vendor Card (Only for customers) */}
          {profile && profile.role === 'customer' && (
            <div className="bg-gradient-to-br from-purple-50 to-pink-50 rounded-xl sm:rounded-2xl shadow-lg p-6 sm:p-8 mb-6 sm:mb-8 border-2 border-purple-200">
              <div className="flex flex-col md:flex-row items-start md:items-center gap-6">
                <div className="flex-shrink-0">
                  <div className="w-16 h-16 sm:w-20 sm:h-20 bg-gradient-to-br from-primary to-secondary rounded-full flex items-center justify-center text-3xl sm:text-4xl">
                    🏪
                  </div>
                </div>
                <div className="flex-1">
                  <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2">
                    Ready to Sell on Nepify?
                  </h3>
                  <p className="text-sm sm:text-base text-gray-700 mb-4">
                    Start your own shop and reach thousands of customers! Join our community of successful vendors.
                  </p>
                  <ul className="text-xs sm:text-sm text-gray-600 space-y-1 mb-4">
                    <li className="flex items-center gap-2">
                      <span className="text-green-500">✓</span>
                      <span>Create your shop in minutes</span>
                    </li>
                    <li className="flex items-center gap-2">
                      <span className="text-green-500">✓</span>
                      <span>List unlimited products</span>
                    </li>
                    <li className="flex items-center gap-2">
                      <span className="text-green-500">✓</span>
                      <span>Manage orders easily</span>
                    </li>
                  </ul>
                  <button
                    onClick={() => setShowVendorModal(true)}
                    className="bg-gradient-to-r from-primary to-secondary text-white px-6 py-3 rounded-lg font-medium hover:shadow-lg transition-all inline-flex items-center gap-2"
                  >
                    <span>🚀</span>
                    <span>Become a Vendor</span>
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Vendor Dashboard (Only for vendors) */}
          {profile && profile.role === 'vendor' && (
            <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl sm:rounded-2xl shadow-lg p-6 sm:p-8 mb-6 sm:mb-8 border-2 border-green-200">
              {shop ? (
                // Has shop - show shop management
                <div>
                  <div className="flex items-center gap-3 mb-4">
                    <span className="text-3xl">🏪</span>
                    <h3 className="text-xl sm:text-2xl font-bold text-gray-900">Your Shop</h3>
                  </div>
                  <div className="bg-white rounded-lg p-6 mb-4">
                    <div className="flex items-start justify-between mb-4">
                      <div>
                        <h4 className="text-2xl font-bold text-gray-900 mb-2">{shop.name}</h4>
                        <p className="text-gray-600 mb-3">{shop.description}</p>
                        <div className="flex items-center gap-3">
                          <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                            shop.is_active 
                              ? 'bg-green-100 text-green-800' 
                              : 'bg-red-100 text-red-800'
                          }`}>
                            {shop.is_active ? 'Active' : 'Inactive'}
                          </span>
                          <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                            shop.is_verified 
                              ? 'bg-blue-100 text-blue-800' 
                              : 'bg-yellow-100 text-yellow-800'
                          }`}>
                            {shop.is_verified ? 'Verified' : 'Pending Verification'}
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="text-sm text-gray-600">
                      <p><span className="font-semibold">Shop URL:</span> nepify.com/shop/{shop.slug}</p>
                      {shop.phone && <p><span className="font-semibold">Phone:</span> {shop.phone}</p>}
                      {shop.email && <p><span className="font-semibold">Email:</span> {shop.email}</p>}
                    </div>
                  </div>
                  <div className="grid sm:grid-cols-3 gap-4">
                    <button className="bg-gradient-to-r from-primary to-secondary text-white px-6 py-3 rounded-lg font-medium hover:shadow-lg transition-all">
                      ➕ Add Product
                    </button>
                    <button className="bg-white border-2 border-gray-300 text-gray-700 px-6 py-3 rounded-lg font-medium hover:bg-gray-50 transition-all">
                      📦 Manage Products
                    </button>
                    <button className="bg-white border-2 border-gray-300 text-gray-700 px-6 py-3 rounded-lg font-medium hover:bg-gray-50 transition-all">
                      ⚙️ Shop Settings
                    </button>
                  </div>
                </div>
              ) : (
                // No shop - show create shop
                <div>
                  <div className="flex items-center gap-3 mb-4">
                    <span className="text-3xl">✅</span>
                    <h3 className="text-xl sm:text-2xl font-bold text-gray-900">Vendor Account Active</h3>
                  </div>
                  <p className="text-sm sm:text-base text-gray-700 mb-6">
                    Your vendor account is ready! Create your shop to start listing products and reaching customers.
                  </p>
                  <div className="grid sm:grid-cols-2 gap-4">
                    <button 
                      onClick={() => navigate('/create-shop')}
                      className="bg-gradient-to-r from-primary to-secondary text-white px-6 py-3 rounded-lg font-medium hover:shadow-lg transition-all"
                    >
                      🏪 Create Your Shop
                    </button>
                    <button className="bg-white border-2 border-gray-300 text-gray-700 px-6 py-3 rounded-lg font-medium hover:bg-gray-50 transition-all">
                      📖 Read Vendor Guide
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}

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

          {/* JWT Token for Testing */}
          <div className="bg-gradient-to-br from-purple-50 to-blue-50 rounded-xl sm:rounded-2xl shadow-lg p-4 sm:p-6 md:p-8 mb-6 sm:mb-8 border-2 border-purple-200">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg sm:text-xl md:text-2xl font-bold text-gray-900">🎫 JWT Token (For Testing)</h3>
              <button
                onClick={loadJwtToken}
                className="bg-purple-600 hover:bg-purple-700 text-white px-3 sm:px-4 py-2 rounded-lg text-xs sm:text-sm font-medium transition-colors"
              >
                🔄 Refresh
              </button>
            </div>
            
            {jwtToken ? (
              <div className="space-y-3">
                <div className="bg-white p-3 sm:p-4 rounded-lg border border-purple-200">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs sm:text-sm font-semibold text-gray-700">Bearer Token:</span>
                    <button
                      onClick={copyTokenToClipboard}
                      className={`px-3 py-1 rounded-md text-xs font-medium transition-colors ${
                        tokenCopied 
                          ? 'bg-green-500 text-white' 
                          : 'bg-purple-600 hover:bg-purple-700 text-white'
                      }`}
                    >
                      {tokenCopied ? '✓ Copied!' : '📋 Copy'}
                    </button>
                  </div>
                  <div className="bg-gray-900 p-2 sm:p-3 rounded overflow-x-auto">
                    <code className="text-xs text-green-400 font-mono break-all whitespace-pre-wrap">
                      {jwtToken}
                    </code>
                  </div>
                </div>
                
                <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3 sm:p-4">
                  <h4 className="text-sm font-semibold text-yellow-800 mb-2">📝 How to use in Insomnia/Postman:</h4>
                  <ol className="text-xs sm:text-sm text-yellow-700 space-y-1 list-decimal list-inside">
                    <li>Click "Copy" button above</li>
                    <li>In Insomnia/Postman, go to Headers tab</li>
                    <li>Add header: <code className="bg-yellow-100 px-1 rounded">Authorization</code></li>
                    <li>Value: <code className="bg-yellow-100 px-1 rounded">Bearer &lt;paste-token-here&gt;</code></li>
                  </ol>
                </div>

                <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 sm:p-4">
                  <h4 className="text-sm font-semibold text-blue-800 mb-2">🧪 Test Endpoints:</h4>
                  <div className="text-xs sm:text-sm text-blue-700 space-y-1">
                    <div><code className="bg-blue-100 px-2 py-0.5 rounded">POST http://localhost:8080/api/v1/shops</code></div>
                    <div><code className="bg-blue-100 px-2 py-0.5 rounded">GET http://localhost:8080/api/v1/my-shop</code></div>
                    <div><code className="bg-blue-100 px-2 py-0.5 rounded">GET http://localhost:8080/api/v1/users/profile</code></div>
                  </div>
                </div>

                <div className="bg-red-50 border border-red-200 rounded-lg p-3 sm:p-4">
                  <h4 className="text-sm font-semibold text-red-800 mb-1">⚠️ Security Warning:</h4>
                  <p className="text-xs sm:text-sm text-red-700">
                    This token is valid for ~1 hour. Never share it publicly or commit it to Git. 
                    This feature is for development testing only.
                  </p>
                </div>
              </div>
            ) : (
              <div className="text-center py-6 text-gray-600">
                <p className="text-sm">Loading token...</p>
              </div>
            )}
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

      {/* Become Vendor Modal */}
      <BecomeVendorModal 
        isOpen={showVendorModal}
        onClose={() => setShowVendorModal(false)}
        onSuccess={handleVendorSuccess}
      />
    </div>
  );
}
