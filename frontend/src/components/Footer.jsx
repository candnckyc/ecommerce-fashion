import React from 'react';

const Footer = () => {
  return (
    <footer className="bg-black text-white mt-20">
      <div className="container mx-auto px-8 py-16">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12">
          {/* Brand */}
          <div>
            <h3 className="text-xl font-light tracking-widest uppercase mb-6">Fashion</h3>
            <p className="text-sm text-gray-400 leading-relaxed">
              Minimalist fashion for the modern individual.
            </p>
          </div>

          {/* Shop */}
          <div>
            <h4 className="text-sm uppercase tracking-wider mb-6">Shop</h4>
            <ul className="space-y-3 text-sm text-gray-400">
              <li><a href="/products" className="hover:text-white transition">All Products</a></li>
              <li><a href="/products?category=4" className="hover:text-white transition">Women</a></li>
              <li><a href="/products?category=5" className="hover:text-white transition">Men</a></li>
            </ul>
          </div>

          {/* Info */}
          <div>
            <h4 className="text-sm uppercase tracking-wider mb-6">Information</h4>
            <ul className="space-y-3 text-sm text-gray-400">
              <li><a href="#" className="hover:text-white transition">About</a></li>
              <li><a href="#" className="hover:text-white transition">Contact</a></li>
              <li><a href="#" className="hover:text-white transition">Shipping</a></li>
            </ul>
          </div>

          {/* Connect */}
          <div>
            <h4 className="text-sm uppercase tracking-wider mb-6">Connect</h4>
            <ul className="space-y-3 text-sm text-gray-400">
              <li><a href="#" className="hover:text-white transition">Instagram</a></li>
              <li><a href="#" className="hover:text-white transition">Twitter</a></li>
              <li><a href="#" className="hover:text-white transition">Pinterest</a></li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 mt-12 pt-8 text-center">
          <p className="text-sm text-gray-500 uppercase tracking-wider">
            © 2025 Fashion. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
