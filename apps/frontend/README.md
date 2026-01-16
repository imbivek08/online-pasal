# Nepify Frontend

Modern React + TypeScript frontend with Clerk authentication integration.

## 🚀 Quick Start

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure Environment Variables

Copy `.env.example` to `.env` and add your Clerk Publishable Key:

```bash
cp .env.example .env
```

Edit `.env`:
```env
VITE_CLERK_PUBLISHABLE_KEY=pk_test_your_key_here
VITE_API_BASE_URL=http://localhost:8080
```

**Get your Clerk key:**
1. Go to [Clerk Dashboard](https://dashboard.clerk.com/)
2. Select your application
3. Go to **API Keys**
4. Copy the **Publishable Key** (starts with `pk_test_`)

### 3. Start Development Server

```bash
npm run dev
```

The app will open at [http://localhost:5173](http://localhost:5173)

## 📁 Project Structure

```
src/
├── lib/
│   └── api.ts              # API client with Clerk auth
├── pages/
│   ├── LandingPage.tsx     # Public landing page
│   ├── SignInPage.tsx      # Clerk sign-in page
│   ├── SignUpPage.tsx      # Clerk sign-up page
│   └── DashboardPage.tsx   # Protected dashboard
├── App.tsx                 # Main app with routing
└── main.tsx               # Entry point with ClerkProvider
```

## 🔐 Authentication Flow

1. **Landing Page** (`/`) - Public welcome page
2. **Sign Up** (`/sign-up`) - Create new account via Clerk
3. **Sign In** (`/sign-in`) - Authenticate via Clerk
4. **Dashboard** (`/dashboard`) - Protected page (requires auth)

## 🔌 API Integration

The app connects to your Go backend API at `http://localhost:8080` by default.

### Using the API Client

```tsx
import { useApi } from './lib/api';

function MyComponent() {
  const api = useApi(); // Automatically includes auth token
  
  const loadData = async () => {
    const response = await api.getProfile();
    console.log(response.data);
  };
  
  return <button onClick={loadData}>Load Profile</button>;
}
```

## 🛠️ Development

### Build for Production
```bash
npm run build
```

### Preview Production Build
```bash
npm run preview
```

## 🔧 Troubleshooting

### "Missing Clerk Publishable Key" Error
Make sure `.env` exists with `VITE_CLERK_PUBLISHABLE_KEY=pk_test_...`

### CORS Errors
Add CORS middleware to your Go backend for `http://localhost:5173`

### "User not found" in Dashboard
Set up Clerk webhooks to sync users (see backend's CLERK_SETUP.md)

## 📖 Learn More

- [Clerk Documentation](https://clerk.com/docs)
- [Clerk React SDK](https://clerk.com/docs/references/react/overview)
- [Vite Documentation](https://vitejs.dev/)

Happy coding! 🎉

