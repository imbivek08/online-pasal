import { useState } from 'react';
import { useApi, type BecomeVendorRequest } from '../lib/api';

interface BecomeVendorModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export default function BecomeVendorModal({ isOpen, onClose, onSuccess }: BecomeVendorModalProps) {
  const api = useApi();
  const [formData, setFormData] = useState<BecomeVendorRequest>({
    business_name: '',
    phone: '',
    business_description: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      const response = await api.becomeVendor(formData);
      
      if (response.success) {
        setSuccess(true);
        setTimeout(() => {
          onSuccess();
          onClose();
          // Reset form
          setFormData({ business_name: '', phone: '', business_description: '' });
          setSuccess(false);
        }, 2000);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to become vendor');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      onClose();
      setError(null);
      setSuccess(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={handleClose}
      />
      
      {/* Modal */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div className="relative bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 sm:p-8">
          {/* Close button */}
          <button
            onClick={handleClose}
            disabled={isSubmitting}
            className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition-colors disabled:opacity-50"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          {success ? (
            // Success state
            <div className="text-center py-8">
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100 mb-4">
                <svg className="h-10 w-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-2">Congratulations! 🎉</h3>
              <p className="text-gray-600 mb-4">
                You are now a vendor! You can create your shop and start selling.
              </p>
              <div className="animate-pulse text-sm text-gray-500">
                Redirecting...
              </div>
            </div>
          ) : (
            // Form state
            <>
              <div className="mb-6">
                <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
                  🏪 Become a Vendor
                </h2>
                <p className="text-sm sm:text-base text-gray-600">
                  Fill in your business information to start selling on Nepify
                </p>
              </div>

              {error && (
                <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-sm text-red-800">⚠️ {error}</p>
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-4">
                {/* Business Name */}
                <div>
                  <label htmlFor="business_name" className="block text-sm font-medium text-gray-700 mb-1">
                    Business Name <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    id="business_name"
                    required
                    minLength={3}
                    maxLength={100}
                    value={formData.business_name}
                    onChange={(e) => setFormData({ ...formData, business_name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    placeholder="e.g., Tech Haven Electronics"
                    disabled={isSubmitting}
                  />
                  <p className="text-xs text-gray-500 mt-1">This will be used for your shop</p>
                </div>

                {/* Phone */}
                <div>
                  <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-1">
                    Phone Number <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="tel"
                    id="phone"
                    required
                    minLength={10}
                    maxLength={20}
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    placeholder="+977-9801234567"
                    disabled={isSubmitting}
                  />
                  <p className="text-xs text-gray-500 mt-1">For customer support and verification</p>
                </div>

                {/* Business Description */}
                <div>
                  <label htmlFor="business_description" className="block text-sm font-medium text-gray-700 mb-1">
                    Business Description <span className="text-gray-400">(Optional)</span>
                  </label>
                  <textarea
                    id="business_description"
                    rows={3}
                    maxLength={500}
                    value={formData.business_description}
                    onChange={(e) => setFormData({ ...formData, business_description: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
                    placeholder="Tell us about your business..."
                    disabled={isSubmitting}
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    {formData.business_description?.length || 0}/500 characters
                  </p>
                </div>

                {/* Agreement */}
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                  <p className="text-xs text-blue-800">
                    ✓ By becoming a vendor, you agree to our Terms of Service and Seller Policy.
                    You'll be able to create your shop and list products immediately.
                  </p>
                </div>

                {/* Buttons */}
                <div className="flex gap-3 pt-2">
                  <button
                    type="button"
                    onClick={handleClose}
                    disabled={isSubmitting}
                    className="flex-1 px-6 py-3 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isSubmitting || !formData.business_name || !formData.phone}
                    className="flex-1 px-6 py-3 bg-gradient-to-r from-primary to-secondary text-white rounded-lg font-medium hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isSubmitting ? (
                      <span className="flex items-center justify-center">
                        <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        Processing...
                      </span>
                    ) : (
                      '🚀 Become Vendor'
                    )}
                  </button>
                </div>
              </form>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
