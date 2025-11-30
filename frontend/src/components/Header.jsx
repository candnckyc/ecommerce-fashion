import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

const Header = () => {
  const { user, logout } = useAuth();
  const { cart } = useCart();

  return (
    <header className="bg-white border-b border-black">
      <div className="container mx-auto px-8">
        <div className="flex items-center justify-between h-20">
          {/* Logo - Minimalist */}
          <Link to="/" className="text-2xl font-light tracking-widest uppercase">
            FASHION
          </Link>

          {/* Navigation - Clean & Spaced */}
          <nav className="hidden md:flex space-x-12">
            <Link to="/products" className="text-sm uppercase tracking-wider hover:opacity-50 transition">
              All
            </Link>
            <Link to="/products?category=1" className="text-sm uppercase tracking-wider hover:opacity-50 transition">
              Women
            </Link>
            <Link to="/products?category=2" className="text-sm uppercase tracking-wider hover:opacity-50 transition">
              Men
            </Link>
          </nav>

          {/* Right side - Minimal icons */}
          <div className="flex items-center space-x-8">
            {user ? (
              <>
                <Link to="/cart" className="relative hover:opacity-50 transition">
                  <span className="text-sm uppercase tracking-wider">Cart</span>
                  {cart.total_items > 0 && (
                    <span className="ml-2 text-xs">({cart.total_items})</span>
                  )}
                </Link>
                <button onClick={logout} className="text-sm uppercase tracking-wider hover:opacity-50 transition">
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm uppercase tracking-wider hover:opacity-50 transition">
                  Login
                </Link>
                <Link to="/register" className="text-sm uppercase tracking-wider hover:opacity-50 transition">
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;