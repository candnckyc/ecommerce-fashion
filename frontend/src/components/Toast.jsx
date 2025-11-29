import React, { useEffect } from 'react';

const Toast = ({ message, type = 'success', onClose }) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, 3000);

    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div className="fixed top-8 right-8 z-50 animate-slide-in">
      <div className={`border-2 bg-white p-6 min-w-[300px] shadow-lg ${
        type === 'error' ? 'border-red-600' : 'border-black'
      }`}>
        <div className="flex items-start justify-between gap-4">
          <p className="text-sm flex-1">{message}</p>
          <button
            onClick={onClose}
            className="text-xl leading-none hover:opacity-50 transition"
          >
            ×
          </button>
        </div>
      </div>
    </div>
  );
};

export default Toast;