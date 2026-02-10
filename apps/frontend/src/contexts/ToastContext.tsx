import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-react';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

interface ToastContextType {
  toast: {
    success: (message: string) => void;
    error: (message: string) => void;
    warning: (message: string) => void;
    info: (message: string) => void;
  };
}

const ToastContext = createContext<ToastContextType | null>(null);

let toastId = 0;

const TOAST_DURATION = 4000;

const toastStyles: Record<ToastType, { bg: string; icon: typeof CheckCircle; border: string }> = {
  success: { bg: 'bg-green-50', icon: CheckCircle, border: 'border-green-400' },
  error: { bg: 'bg-red-50', icon: XCircle, border: 'border-red-400' },
  warning: { bg: 'bg-yellow-50', icon: AlertTriangle, border: 'border-yellow-400' },
  info: { bg: 'bg-blue-50', icon: Info, border: 'border-blue-400' },
};

const iconColors: Record<ToastType, string> = {
  success: 'text-green-500',
  error: 'text-red-500',
  warning: 'text-yellow-500',
  info: 'text-blue-500',
};

const textColors: Record<ToastType, string> = {
  success: 'text-green-800',
  error: 'text-red-800',
  warning: 'text-yellow-800',
  info: 'text-blue-800',
};

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback(
    (message: string, type: ToastType) => {
      const id = ++toastId;
      setToasts((prev) => [...prev, { id, message, type }]);
      setTimeout(() => removeToast(id), TOAST_DURATION);
    },
    [removeToast]
  );

  const toast = {
    success: (message: string) => addToast(message, 'success'),
    error: (message: string) => addToast(message, 'error'),
    warning: (message: string) => addToast(message, 'warning'),
    info: (message: string) => addToast(message, 'info'),
  };

  return (
    <ToastContext.Provider value={{ toast }}>
      {children}

      {/* Toast container — fixed top-right */}
      <div className="fixed top-4 right-4 z-[9999] flex flex-col gap-3 pointer-events-none">
        {toasts.map((t) => {
          const style = toastStyles[t.type];
          const Icon = style.icon;

          return (
            <div
              key={t.id}
              className={`pointer-events-auto flex items-start gap-3 min-w-[320px] max-w-[420px] px-4 py-3 rounded-lg border ${style.bg} ${style.border} shadow-lg animate-slide-in`}
            >
              <Icon className={`w-5 h-5 mt-0.5 shrink-0 ${iconColors[t.type]}`} />
              <p className={`text-sm font-medium flex-1 ${textColors[t.type]}`}>{t.message}</p>
              <button
                onClick={() => removeToast(t.id)}
                className="shrink-0 text-gray-400 hover:text-gray-600 transition"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          );
        })}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context.toast;
}
