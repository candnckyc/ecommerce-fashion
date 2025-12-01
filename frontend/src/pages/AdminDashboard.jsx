import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || '185.26.144.214:4501';

const AdminDashboard = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Redirect if not admin
    if (user && user.role !== 'admin') {
      navigate('/');
      return;
    }

    // Fetch stats
    const token = localStorage.getItem('token');
    axios.get(`${API_URL}/admin/stats`, {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then(response => {
        setStats(response.data.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load stats:', error);
        setLoading(false);
      });
  }, [user, navigate]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-xl uppercase tracking-wider">Loading...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-12">
      <h1 className="heading-lg mb-12">Admin Dashboard</h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-12">
        <div className="border border-black p-6">
          <div className="text-sm uppercase tracking-wider text-gray-600 mb-2">
            Total Products
          </div>
          <div className="text-3xl font-light">{stats?.total_products || 0}</div>
        </div>
        
        <div className="border border-black p-6">
          <div className="text-sm uppercase tracking-wider text-gray-600 mb-2">
            Total Orders
          </div>
          <div className="text-3xl font-light">{stats?.total_orders || 0}</div>
        </div>
        
        <div className="border border-black p-6">
          <div className="text-sm uppercase tracking-wider text-gray-600 mb-2">
            Pending Orders
          </div>
          <div className="text-3xl font-light">{stats?.pending_orders || 0}</div>
        </div>
        
        <div className="border border-black p-6">
          <div className="text-sm uppercase tracking-wider text-gray-600 mb-2">
            Total Revenue
          </div>
          <div className="text-3xl font-light">
            {(stats?.total_revenue || 0).toFixed(2)} TL
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <button
          onClick={() => navigate('/admin/products')}
          className="border border-black p-8 hover:bg-black hover:text-white transition-colors"
        >
          <div className="text-lg uppercase tracking-wider">Manage Products</div>
          <div className="text-sm mt-2 opacity-70">Add, edit, or remove products</div>
        </button>
        
        <button
          onClick={() => navigate('/admin/orders')}
          className="border border-black p-8 hover:bg-black hover:text-white transition-colors"
        >
          <div className="text-lg uppercase tracking-wider">Manage Orders</div>
          <div className="text-sm mt-2 opacity-70">View and update order status</div>
        </button>
        
        <button
          onClick={() => navigate('/admin/customers')}
          className="border border-black p-8 hover:bg-black hover:text-white transition-colors"
        >
          <div className="text-lg uppercase tracking-wider">View Customers</div>
          <div className="text-sm mt-2 opacity-70">Customer list and details</div>
        </button>
      </div>
    </div>
  );
};

export default AdminDashboard;