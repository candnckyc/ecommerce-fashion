import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import SearchBar from './SearchBar';

const Header = () => {
  const { user, logout } = useAuth();
  const { cart } = useCart();
  const [activeDropdown, setActiveDropdown] = useState(null);

  const categories = {
    women: [
      { id: 4, name: "Dresses" },
      { id: 5, name: "Tops" },
      { id: 6, name: "Jeans" },
      { id: 7, name: "Outerwear" }
    ],
    men: [
      { id: 8, name: "Shirts" },
      { id: 9, name: "T-Shirts" },
      { id: 10, name: "Pants" },
      { id: 11, name: "Jackets" }
    ],
    kids: [
      { id: 12, name: "Girls" },
      { id: 13, name: "Boys" }
    ]
  };

  return (
    <header className="bg-white border-b border-black">
      <div className="container mx-auto px-8">
        <div className="flex items-center justify-between h-20">
          {/* Logo - Minimalist */}
          <Link to="/" className="text-2xl font-light tracking-widest uppercase">
            FASHION
          </Link>

          {/* Navigation - Clean & Spaced with Dropdowns */}
          <nav className="hidden md:flex space-x-12 relative">
            <Link 
              to="/products" 
              className="text-sm uppercase tracking-wider hover:opacity-50 transition inline-block py-2"
            >
              All
            </Link>
            
            {/* Women Dropdown */}
            <div 
              className="relative"
              onMouseEnter={() => setActiveDropdown('women')}
              onMouseLeave={() => setActiveDropdown(null)}
            >
              <Link 
                to="/products?category=1" 
                className="text-sm uppercase tracking-wider hover:opacity-50 transition inline-block py-2"
              >
                Women
              </Link>
              
              {activeDropdown === 'women' && (
                <div className="absolute top-full left-0 w-48 bg-white border border-black shadow-lg z-50">
                  {categories.women.map(cat => (
                    <Link
                      key={cat.id}
                      to={`/products?category=${cat.id}`}
                      className="block px-4 py-2 text-xs uppercase tracking-wider hover:bg-black hover:text-white border-b border-gray-200 last:border-b-0"
                    >
                      {cat.name}
                    </Link>
                  ))}
                </div>
              )}
            </div>

            {/* Men Dropdown */}
            <div 
              className="relative"
              onMouseEnter={() => setActiveDropdown('men')}
              onMouseLeave={() => setActiveDropdown(null)}
            >
              <Link 
                to="/products?category=2" 
                className="text-sm uppercase tracking-wider hover:opacity-50 transition inline-block py-2"
              >
                Men
              </Link>
              
              {activeDropdown === 'men' && (
                <div className="absolute top-full left-0 w-48 bg-white border border-black shadow-lg z-50">
                  {categories.men.map(cat => (
                    <Link
                      key={cat.id}
                      to={`/products?category=${cat.id}`}
                      className="block px-4 py-2 text-xs uppercase tracking-wider hover:bg-black hover:text-white border-b border-gray-200 last:border-b-0"
                    >
                      {cat.name}
                    </Link>
                  ))}
                </div>
              )}
            </div>

            {/* Kids Dropdown */}
            <div 
              className="relative"
              onMouseEnter={() => setActiveDropdown('kids')}
              onMouseLeave={() => setActiveDropdown(null)}
            >
              <Link 
                to="/products?category=3" 
                className="text-sm uppercase tracking-wider hover:opacity-50 transition inline-block py-2"
              >
                Kids
              </Link>
              
              {activeDropdown === 'kids' && (
                <div className="absolute top-full left-0 w-48 bg-white border border-black shadow-lg z-50">
                  {categories.kids.map(cat => (
                    <Link
                      key={cat.id}
                      to={`/products?category=${cat.id}`}
                      className="block px-4 py-2 text-xs uppercase tracking-wider hover:bg-black hover:text-white border-b border-gray-200 last:border-b-0"
                    >
                      {cat.name}
                    </Link>
                  ))}
                </div>
              )}
            </div>
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
                {user.role === 'admin' && (
                  <Link to="/admin" className="text-sm uppercase tracking-wider hover:opacity-50 transition">
                    Admin
                  </Link>
                )}
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

        {/* Search Bar - Below main header */}
        <div className="pb-4">
          <SearchBar />
        </div>
      </div>
    </header>
  );
};

export default Header;