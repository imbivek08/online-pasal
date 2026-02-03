import { SignedIn, SignedOut, UserButton, useUser } from '@clerk/clerk-react';
import { Link, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import CartIcon from './CartIcon';

export default function Navbar() {
  const { user } = useUser();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/products?search=${encodeURIComponent(searchQuery)}`);
      setMobileMenuOpen(false);
    }
  };

  return (
    <nav className="bg-white shadow-md fixed top-0 left-0 right-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center space-x-2 hover:opacity-80 transition-opacity flex-shrink-0">
            <span className="text-2xl">🚀</span>
            <span className="text-xl sm:text-2xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              Nepify
            </span>
          </Link>

          {/* Desktop Search Bar */}
          <form onSubmit={handleSearch} className="flex-1 max-w-md mx-8 hidden md:block">
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                type="text"
                placeholder="Search products..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-2 focus:ring-primary focus:border-transparent text-gray-900"
              />
            </div>
          </form>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-4 lg:space-x-6">
            <Link to="/products" className="text-gray-700 hover:text-primary font-medium transition-colors text-sm lg:text-base">
              Products
            </Link>
            
            <SignedIn>
              <Link to="/dashboard" className="text-gray-700 hover:text-primary font-medium transition-colors text-sm lg:text-base">
                Dashboard
              </Link>
              <CartIcon />
            </SignedIn>

            <SignedOut>
              <Link 
                to="/sign-in" 
                className="text-gray-700 hover:text-primary font-medium transition-colors text-sm lg:text-base"
              >
                Sign In
              </Link>
              <Link 
                to="/sign-up" 
                className="bg-gradient-to-r from-primary to-secondary text-white px-3 lg:px-4 py-2 rounded-lg font-medium hover:shadow-lg transition-all text-sm lg:text-base"
              >
                Sign Up
              </Link>
            </SignedOut>
            
            <SignedIn>
              <div className="flex items-center space-x-3">
                <span className="text-gray-700 font-medium hidden lg:block text-sm">
                  Hi, {user?.firstName || user?.emailAddresses[0].emailAddress.split('@')[0]}
                </span>
                <UserButton 
                  afterSignOutUrl="/"
                  appearance={{
                    elements: {
                      avatarBox: "w-10 h-10"
                    }
                  }}
                />
              </div>
            </SignedIn>
          </div>

          {/* Mobile Menu Button */}
          <div className="flex md:hidden items-center space-x-2">
            <SignedIn>
              <UserButton 
                afterSignOutUrl="/"
                appearance={{
                  elements: {
                    avatarBox: "w-9 h-9"
                  }
                }}
              />
            </SignedIn>
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="text-gray-700 hover:text-primary p-2 rounded-lg hover:bg-gray-100 transition-colors"
              aria-label="Toggle menu"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                {mobileMenuOpen ? (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                ) : (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                )}
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <div className="md:hidden border-t border-gray-200 py-4 space-y-4">
            {/* Mobile Search */}
            <form onSubmit={handleSearch} className="px-2">
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
                <input
                  type="text"
                  placeholder="Search products..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-gray-900"
                />
              </div>
            </form>

            {/* Mobile Links */}
            <div className="flex flex-col space-y-2 px-2">
              <Link 
                to="/products" 
                onClick={() => setMobileMenuOpen(false)}
                className="text-gray-700 hover:text-primary hover:bg-gray-50 font-medium py-3 px-3 rounded-lg transition-colors"
              >
                Products
              </Link>
              
              <SignedIn>
                <Link 
                  to="/dashboard" 
                  onClick={() => setMobileMenuOpen(false)}
                  className="text-gray-700 hover:text-primary hover:bg-gray-50 font-medium py-3 px-3 rounded-lg transition-colors"
                >
                  Dashboard
                <Link 
                  to="/cart" 
                  onClick={() => setMobileMenuOpen(false)}
                  className="text-gray-700 hover:text-primary hover:bg-gray-50 font-medium py-3 px-3 rounded-lg transition-colors"
                >
                  Cart
                </Link>
                </Link>
              </SignedIn>

              <SignedOut>
                <Link 
                  to="/sign-in" 
                  onClick={() => setMobileMenuOpen(false)}
                  className="text-gray-700 hover:text-primary hover:bg-gray-50 font-medium py-3 px-3 rounded-lg transition-colors"
                >
                  Sign In
                </Link>
                <Link 
                  to="/sign-up" 
                  onClick={() => setMobileMenuOpen(false)}
                  className="bg-gradient-to-r from-primary to-secondary text-white text-center py-3 px-3 rounded-lg font-medium hover:shadow-lg transition-all"
                >
                  Sign Up
                </Link>
              </SignedOut>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
