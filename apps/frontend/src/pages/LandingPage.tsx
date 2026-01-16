import { Link } from 'react-router-dom';
import Navbar from '../components/Navbar';

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-primary via-secondary to-purple-900">
      <Navbar />
      
      {/* Hero Section */}
      <div className="pt-24 sm:pt-32 pb-12 sm:pb-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto text-center">
          <h1 className="text-3xl sm:text-5xl md:text-6xl lg:text-7xl font-bold text-white mb-4 sm:mb-6 leading-tight">
            🚀 Welcome to <span className="text-accent">Nepify</span>
          </h1>
          <p className="text-base sm:text-xl md:text-2xl text-white/90 mb-8 sm:mb-12 max-w-3xl mx-auto px-2">
            Your modern e-commerce platform for the future. 
            Buy and sell with confidence.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-stretch sm:items-center mb-12 sm:mb-16 px-4">
            <Link 
              to="/sign-up" 
              className="bg-white text-primary px-6 sm:px-8 py-3 sm:py-4 rounded-full text-base sm:text-lg font-bold shadow-2xl hover:shadow-accent/50 hover:scale-105 transition-all duration-300 text-center"
            >
              Get Started Free
            </Link>
            <Link 
              to="/products" 
              className="bg-transparent text-white px-6 sm:px-8 py-3 sm:py-4 rounded-full text-base sm:text-lg font-bold border-2 border-white hover:bg-white/10 transition-all duration-300 text-center"
            >
              Browse Products
            </Link>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="bg-white py-12 sm:py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold text-center text-gray-900 mb-8 sm:mb-16">
            Why Choose Nepify?
          </h2>
          
          <div className="grid sm:grid-cols-2 md:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
            {/* Feature 1 */}
            <div className="bg-gradient-to-br from-primary/10 to-secondary/10 p-6 sm:p-8 rounded-2xl hover:shadow-xl transition-all duration-300 hover:-translate-y-2">
              <div className="text-4xl sm:text-6xl mb-3 sm:mb-4">🔐</div>
              <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2 sm:mb-3">
                Secure Authentication
              </h3>
              <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
                Enterprise-grade security powered by Clerk. Your data is always protected with industry-leading encryption.
              </p>
            </div>

            {/* Feature 2 */}
            <div className="bg-gradient-to-br from-primary/10 to-secondary/10 p-6 sm:p-8 rounded-2xl hover:shadow-xl transition-all duration-300 hover:-translate-y-2">
              <div className="text-4xl sm:text-6xl mb-3 sm:mb-4">⚡</div>
              <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2 sm:mb-3">
                Lightning Fast
              </h3>
              <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
                Built with Go and React for optimal performance. Experience blazing fast load times and smooth interactions.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="bg-gradient-to-br from-primary/10 to-secondary/10 p-6 sm:p-8 rounded-2xl hover:shadow-xl transition-all duration-300 hover:-translate-y-2">
              <div className="text-4xl sm:text-6xl mb-3 sm:mb-4">🎨</div>
              <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2 sm:mb-3">
                Modern Design
              </h3>
              <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
                Beautiful, intuitive interface designed for the best user experience. Clean and professional aesthetics.
              </p>
            </div>

            {/* Feature 4 */}
            <div className="bg-gradient-to-br from-primary/10 to-secondary/10 p-6 sm:p-8 rounded-2xl hover:shadow-xl transition-all duration-300 hover:-translate-y-2">
              <div className="text-4xl sm:text-6xl mb-3 sm:mb-4">💼</div>
              <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2 sm:mb-3">
                Vendor Dashboard
              </h3>
              <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
                Manage your products, track sales, and grow your business with our comprehensive vendor tools.
              </p>
            </div>

            {/* Feature 5 */}
            <div className="bg-gradient-to-br from-primary/10 to-secondary/10 p-6 sm:p-8 rounded-2xl hover:shadow-xl transition-all duration-300 hover:-translate-y-2">
              <div className="text-4xl sm:text-6xl mb-3 sm:mb-4">🛒</div>
              <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2 sm:mb-3">
                Easy Shopping
              </h3>
              <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
                Streamlined checkout process. Find what you need quickly with powerful search and filtering.
              </p>
            </div>

            {/* Feature 6 */}
            <div className="bg-gradient-to-br from-primary/10 to-secondary/10 p-6 sm:p-8 rounded-2xl hover:shadow-xl transition-all duration-300 hover:-translate-y-2">
              <div className="text-4xl sm:text-6xl mb-3 sm:mb-4">📱</div>
              <h3 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2 sm:mb-3">
                Mobile Friendly
              </h3>
              <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
                Fully responsive design works perfectly on all devices. Shop on the go with ease.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="bg-gradient-to-r from-primary to-secondary py-12 sm:py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold text-white mb-4 sm:mb-6">
            Ready to Start Selling?
          </h2>
          <p className="text-base sm:text-lg md:text-xl text-white/90 mb-6 sm:mb-8 px-4">
            Join thousands of vendors already growing their business on Nepify.
          </p>
          <Link 
            to="/sign-up" 
            className="inline-block bg-white text-primary px-8 sm:px-10 py-3 sm:py-4 rounded-full text-base sm:text-lg font-bold shadow-2xl hover:shadow-accent/50 hover:scale-105 transition-all duration-300"
          >
            Create Your Account
          </Link>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-8 sm:py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto text-center">
          <p className="text-sm sm:text-base text-gray-400">
            © 2024 Nepify. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
