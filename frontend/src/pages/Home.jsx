import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { productsAPI } from '../services/api';

const Home = () => {
  const [featuredProducts, setFeaturedProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    productsAPI.getAll({ limit: 4 })
      .then(response => {
        setFeaturedProducts(response.data.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load products:', error);
        setLoading(false);
      });
  }, []);

  return (
    <div>
      {/* Hero Section - Minimalist */}
      <section className="min-h-screen flex items-center justify-center bg-white border-b border-black">
        <div className="text-center px-8">
          <h1 className="heading-xl mb-8">
            Essential<br/>Minimalism
          </h1>
          <p className="text-sm tracking-wider uppercase mb-12 max-w-md mx-auto">
            Timeless pieces for the modern wardrobe
          </p>
          <Link to="/products" className="btn-primary">
            Explore Collection
          </Link>
        </div>
      </section>

      {/* Featured Products - Grid */}
      <section className="container mx-auto px-8 py-20">
        <h2 className="heading-md mb-16 text-center">Featured</h2>
        
        {loading ? (
          <div className="text-center py-20">
            <p className="uppercase tracking-wider text-sm">Loading...</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-px bg-black">
            {featuredProducts.map(product => (
              <Link 
                key={product.id} 
                to={`/products/${product.id}`}
                className="group bg-white"
              >
                <div className="aspect-[3/4] bg-gray-100 overflow-hidden">
                  {product.images && product.images[0] ? (
                    <img 
                      src={product.images[0].image_url} 
                      alt={product.name}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <span className="text-sm uppercase tracking-wider text-gray-400">No Image</span>
                    </div>
                  )}
                </div>
                <div className="p-6">
                  <h3 className="text-sm uppercase tracking-wider mb-2">{product.name}</h3>
                  <p className="text-sm">{product.base_price.toFixed(2)} TL</p>
                </div>
              </Link>
            ))}
          </div>
        )}
      </section>

      {/* Philosophy Section */}
      <section className="bg-black text-white py-32">
        <div className="container mx-auto px-8 max-w-2xl text-center">
          <h2 className="heading-md mb-8">Philosophy</h2>
          <p className="text-sm tracking-wide leading-relaxed opacity-80">
            We believe in the power of simplicity. Each piece is designed to transcend trends,
            focusing on form, function, and timeless elegance. Our collections celebrate the
            beauty of reduction and the strength of minimalism.
          </p>
        </div>
      </section>
    </div>
  );
};

export default Home;
