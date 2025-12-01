import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || '185.26.144.214:4501';

const AdminOrders = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Redirect if not admin
    if (user && user.role !== 'admin') {
      navigate('/');
      return;
    }

    // Fetch all orders (admin endpoint)
    const token = localStorage.getItem('token');
    axios.get(`${API_URL}/admin/orders`, {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then(response => {
        setOrders(response.data.data || []);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load orders:', error);
        setLoading(false);
      });
  }, [user, navigate]);

  const handleUpdateStatus = (orderId, newStatus) => {
    const token = localStorage.getItem('token');
    axios.put(`${API_URL}/admin/orders/${orderId}/status`, 
      { status: newStatus },
      { headers: { Authorization: `Bearer ${token}` }}
    )
      .then(() => {
        // Update local state
        setOrders(orders.map(order => 
          order.id === orderId ? { ...order, status: newStatus } : order
        ));
      })
      .catch(error => {
        console.error('Failed to update status:', error);
      });
  };

  const getStatusColor = (status) => {
    const colors = {
      pending: 'bg-gray-100 text-gray-800',
      confirmed: 'bg-blue-100 text-blue-800',
      processing: 'bg-yellow-100 text-yellow-800',
      shipped: 'bg-purple-100 text-purple-800',
      delivered: 'bg-green-100 text-green-800',
      cancelled: 'bg-red-100 text-red-800',
    };
    return colors[status] || 'bg-gray-100 text-gray-800';
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
        <h1 className="heading-lg">Manage Orders</h1>
        <button
          onClick={() => navigate('/admin')}
          className="btn-secondary"
        >
          Back to Dashboard
        </button>
      </div>

      <div className="mb-6 text-sm uppercase tracking-wider text-gray-600">
        Total Orders: {orders.length}
      </div>

      {/* Orders Table */}
      <div className="border border-black">
        <table className="w-full">
          <thead className="border-b border-black bg-gray-50">
            <tr>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Order #</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Customer</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Total</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Status</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Date</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.length === 0 ? (
              <tr>
                <td colSpan="6" className="px-6 py-8 text-center text-gray-500">
                  No orders found
                </td>
              </tr>
            ) : (
              orders.map((order) => (
                <tr key={order.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium">{order.order_number}</td>
                  <td className="px-6 py-4 text-sm">{order.shipping_full_name}</td>
                  <td className="px-6 py-4 text-sm">{order.total?.toFixed(2)} TL</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs uppercase tracking-wider px-2 py-1 ${getStatusColor(order.status)}`}>
                      {order.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(order.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <select
                      value={order.status}
                      onChange={(e) => handleUpdateStatus(order.id, e.target.value)}
                      className="border border-black px-2 py-1 text-xs uppercase tracking-wider"
                    >
                      <option value="pending">Pending</option>
                      <option value="confirmed">Confirmed</option>
                      <option value="processing">Processing</option>
                      <option value="shipped">Shipped</option>
                      <option value="delivered">Delivered</option>
                      <option value="cancelled">Cancelled</option>
                    </select>
                    <button
                      onClick={() => navigate(`/order-confirmation/${order.id}`)}
                      className="ml-4 text-blue-600 hover:underline"
                    >
                      View
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AdminOrders;