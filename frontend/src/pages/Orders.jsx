import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { orderAPI } from '../services/api';

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    orderAPI.getOrders()
      .then(response => {
        setOrders(response.data.data || []);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load orders:', error);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-sm uppercase tracking-wider">Loading...</p>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="heading-md mb-8">No Orders Yet</h1>
          <Link to="/products" className="btn-primary">
            Start Shopping
          </Link>
        </div>
      </div>
    );
  }

  const getStatusColor = (status) => {
    const colors = {
      pending: 'text-gray-600',
      confirmed: 'text-blue-600',
      processing: 'text-yellow-600',
      shipped: 'text-purple-600',
      delivered: 'text-green-600',
      cancelled: 'text-red-600',
    };
    return colors[status] || 'text-gray-600';
  };

  return (
    <div className="container mx-auto px-8 py-12">
      <h1 className="heading-lg mb-12">Order History</h1>

      <div className="space-y-6">
        {orders.map((order) => (
          <div key={order.id} className="border border-black hover:border-2 transition">
            <div className="p-6">
              {/* Order Header */}
              <div className="flex justify-between items-start mb-4">
                <div>
                  <p className="text-sm font-bold mb-1">Order #{order.order_number}</p>
                  <p className="text-xs text-gray-600">
                    {new Date(order.created_at).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric',
                    })}
                  </p>
                </div>
                <div className="text-right">
                  <p className={`text-sm uppercase tracking-wider ${getStatusColor(order.status)}`}>
                    {order.status}
                  </p>
                  <p className="text-sm mt-1">{order.total.toFixed(2)} TL</p>
                </div>
              </div>

              {/* Shipping Address */}
              <div className="text-xs text-gray-600 mb-4">
                <p>Ship to: {order.shipping_full_name}</p>
                <p>{order.shipping_address_line1}, {order.shipping_city}, {order.shipping_country}</p>
              </div>

              {/* Actions */}
              <div className="flex gap-4">
                <Link
                  to={`/order-confirmation/${order.id}`}
                  className="text-sm uppercase tracking-wider border border-black px-6 py-2 hover:bg-black hover:text-white transition"
                >
                  View Details
                </Link>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Orders;