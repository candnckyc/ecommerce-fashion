import React, { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { productsAPI } from '../services/api';

const ProductList = () => {
  const [products, setProducts] = useState([]);
  const [brands, setBrands] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    // Fetch filters
    Promise.all([
      productsAPI.getBrands(),
      productsAPI.getCategories()
    ]).then(([brandsRes, categoriesRes]) => {
      setBrands(brandsRes.data.data);
      setCategories(categoriesRes.data.data);
    });
  }, []);

  useEffect(() => {
    // Fetch products with filters
    setLoading(true);
    const params = Object.fromEntries(searchParams);
    productsAPI.getAll(params)
      .then(response => {
        setProducts(response.data.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load products:', error);
        setLoading(false);
      });
  }, [searchParams]);

  const handleFilterChange = (key, value) => {
    const newParams = new URLSearchParams(searchParams);
    if (value) {
      newParams.set(key, value);
    } else {
      newParams.delete(key);
    }
    setSearchParams(newParams);
  };

  return (
    <div className="container mx-auto px-8 py-12">
      <h1 className="heading-lg mb-12">Collection</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-12">
        {/* Filters Sidebar */}
        <aside className="md:col-span-1">
          <div className="border border-black p-6">
            <h3 className="text-sm uppercase tracking-wider mb-6">Filter</h3>

            {/* Categories */}
            <div className="mb-8">
              <h4 className="text-xs uppercase tracking-wider mb-4 text-gray-600">Category</h4>
              <select 
                className="w-full border border-black p-2 text-sm"
                value={searchParams.get('category') || ''}
                onChange={(e) => handleFilterChange('category', e.target.value)}
              >
                <option value="">All Categories</option>
                {categories.map(cat => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>

            {/* Brands */}
            <div className="mb-8">
              <h4 className="text-xs uppercase tracking-wider mb-4 text-gray-600">Brand</h4>
              <select 
                className="w-full border border-black p-2 text-sm"
                value={searchParams.get('brand') || ''}
                onChange={(e) => handleFilterChange('brand', e.target.value)}
              >
                <option value="">All Brands</option>
                {brands.map(brand => (
                  <option key={brand.id} value={brand.id}>{brand.name}</option>
                ))}
              </select>
            </div>

            {/* Clear Filters */}
            <button 
              onClick={() => setSearchParams({})}
              className="w-full text-xs uppercase tracking-wider border border-black p-2 hover:bg-black hover:text-white transition"
            >
              Clear All
            </button>
          </div>
        </aside>

        {/* Products Grid */}
        <div className="md:col-span-3">
          {loading ? (
            <div className="text-center py-20">
              <p className="uppercase tracking-wider text-sm">Loading...</p>
            </div>
          ) : products.length === 0 ? (
            <div className="text-center py-20">
              <p className="uppercase tracking-wider text-sm">No products found</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-px bg-black">
              {products.map(product => (
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
        </div>
      </div>
    </div>
  );
};

export default ProductList;
