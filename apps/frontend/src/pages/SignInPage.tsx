import { SignIn } from '@clerk/clerk-react';
import { Link } from 'react-router-dom';

export default function SignInPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-primary via-secondary to-purple-900 flex items-center justify-center p-3 sm:p-4">
      <div className="w-full max-w-md">
        {/* Logo and Back Link */}
        <div className="text-center mb-6 sm:mb-8">
          <Link to="/" className="inline-flex items-center space-x-2 text-white hover:opacity-80 transition-opacity">
            <span className="text-3xl sm:text-4xl">🚀</span>
            <span className="text-2xl sm:text-3xl font-bold">Nepify</span>
          </Link>
          <p className="text-sm sm:text-base text-white/80 mt-2">Welcome back! Sign in to your account</p>
        </div>
        
        {/* Clerk SignIn Component */}
        <div className="flex justify-center">
          <SignIn 
            routing="path" 
            path="/sign-in"
            signUpUrl="/sign-up"
            appearance={{
              elements: {
                rootBox: "w-full",
                card: "shadow-2xl",
              }
            }}
          />
        </div>
      </div>
    </div>
  );
}
