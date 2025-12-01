import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { productsAPI } from '../services/api';

const AdminProducts = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Redirect if not admin
    if (user && user.role !== 'admin') {
      navigate('/');
      return;
    }

    // Fetch all products
    productsAPI.getAll({ limit: 200 })
      .then(response => {
        setProducts(response.data.data || []);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load products:', error);
        setLoading(false);
      });
  }, [user, navigate]);

  const handleToggleActive = (productId) => {
    const token = localStorage.getItem('token');
    const API_URL = process.env.REACT_APP_API_URL || '185.26.144.214:4501';
    
    fetch(`${API_URL}/admin/products/${productId}/toggle`, {
      method: 'PUT',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
      .then(response => response.json())
      .then(() => {
        // Update local state
        setProducts(products.map(p => 
          p.id === productId ? { ...p, is_active: !p.is_active } : p
        ));
      })
      .catch(error => {
        console.error('Failed to toggle product:', error);
        alert('Failed to toggle product status');
      });
  };

  const handleDelete = (productId) => {
    if (window.confirm('Are you sure you want to delete this product?')) {
      // TODO: Implement delete
      console.log('Delete product:', productId);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl uppercase tracking-wider">Loading...</div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-12">
      <div className="flex justify-between items-center mb-12">
        <h1 className="heading-lg">Manage Products</h1>
        <button
          onClick={() => navigate('/admin')}
          className="btn-secondary"
        >
          Back to Dashboard
        </button>
      </div>

      <div className="mb-6 text-sm uppercase tracking-wider text-gray-600">
        Total Products: {products.length}
      </div>

      {/* Products Table */}
      <div className="border border-black">
        <table className="w-full">
          <thead className="border-b border-black bg-gray-50">
            <tr>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">ID</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Name</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Price</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Status</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {products.map((product) => (
              <tr key={product.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm">{product.id}</td>
                <td className="px-6 py-4">
                  <div className="flex items-center">
                    {product.images && product.images[0] && (
                      <img
                        src={product.images[0].image_url}
                        alt={product.name}
                        className="w-12 h-12 object-cover border border-black mr-4"
                      />
                    )}
                    <div>
                      <div className="text-sm font-medium">{product.name}</div>
                      <div className="text-xs text-gray-500">{product.slug}</div>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 text-sm">{product.base_price.toFixed(2)} TL</td>
                <td className="px-6 py-4">
                  <span className={`text-xs uppercase tracking-wider px-2 py-1 ${
                    product.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                  }`}>
                    {product.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm">
                  <button
                    onClick={() => navigate(`/products/${product.id}`)}
                    className="text-blue-600 hover:underline mr-4"
                  >
                    View
                  </button>
                  <button
                    onClick={() => handleToggleActive(product.id)}
                    className="text-yellow-600 hover:underline mr-4"
                  >
                    Toggle
                  </button>
                  <button
                    onClick={() => handleDelete(product.id)}
                    className="text-red-600 hover:underline"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AdminProducts;