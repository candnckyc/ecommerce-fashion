import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { orderAPI } from '../services/api';

const OrderConfirmation = () => {
  const { id } = useParams();
  const [order, setOrder] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    orderAPI.getOrderById(id)
      .then(response => {
        setOrder(response.data.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load order:', error);
        setLoading(false);
      });
  }, [id]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-sm uppercase tracking-wider">Loading...</p>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-sm uppercase tracking-wider">Order not found</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-8 py-12">
      <div className="max-w-3xl mx-auto">
        {/* Success Message */}
        <div className="text-center mb-12">
          <div className="w-16 h-16 border-2 border-black mx-auto mb-6 flex items-center justify-center">
            <span className="text-3xl">✓</span>
          </div>
          <h1 className="heading-lg mb-4">Order Confirmed</h1>
          <p className="text-sm tracking-wide mb-2">Thank you for your order</p>
          <p className="text-sm tracking-wide text-gray-600">Order #{order.order_number}</p>
        </div>

        {/* Order Details */}
        <div className="border border-black p-8 mb-8">
          <h2 className="text-sm uppercase tracking-wider mb-6">Order Details</h2>

          {/* Items */}
          <div className="space-y-4 mb-8">
            {order.items?.map((item) => (
              <div key={item.id} className="flex justify-between text-sm">
                <div>
                  <p className="font-bold">{item.product_name}</p>
                  <p className="text-gray-600">Size: {item.size} | Color: {item.color}</p>
                  <p className="text-gray-600">Quantity: {item.quantity}</p>
                </div>
                <div className="text-right">
                  <p>{item.total_price.toFixed(2)} TL</p>
                </div>
              </div>
            ))}
          </div>

          {/* Totals */}
          <div className="border-t border-black pt-4 space-y-2">
            <div className="flex justify-between text-sm">
              <span>Subtotal</span>
              <span>{order.subtotal.toFixed(2)} TL</span>
            </div>
            <div className="flex justify-between text-sm">
              <span>Shipping</span>
              <span>{order.shipping_cost.toFixed(2)} TL</span>
            </div>
            <div className="flex justify-between text-lg font-bold border-t border-black pt-4">
              <span className="uppercase tracking-wider">Total</span>
              <span>{order.total.toFixed(2)} TL</span>
            </div>
          </div>
        </div>

        {/* Shipping Address */}
        <div className="border border-black p-8 mb-8">
          <h2 className="text-sm uppercase tracking-wider mb-4">Shipping Address</h2>
          <div className="text-sm">
            <p className="font-bold">{order.shipping_full_name}</p>
            <p>{order.shipping_address_line1}</p>
            {order.shipping_address_line2 && <p>{order.shipping_address_line2}</p>}
            <p>{order.shipping_city}, {order.shipping_state} {order.shipping_postal_code}</p>
            <p>{order.shipping_country}</p>
            <p className="mt-2">{order.shipping_phone}</p>
          </div>
        </div>

        {/* Payment Info */}
        <div className="border border-black p-8 mb-8">
          <h2 className="text-sm uppercase tracking-wider mb-4">Payment Information</h2>
          <div className="text-sm">
            <p>Method: <span className="uppercase">{order.payment_method.replace('_', ' ')}</span></p>
            <p>Status: <span className="uppercase">{order.payment_status}</span></p>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-4 justify-center">
          <Link to="/orders" className="btn-secondary">
            View All Orders
          </Link>
          <Link to="/products" className="btn-primary">
            Continue Shopping
          </Link>
        </div>
      </div>
    </div>
  );
};

export default OrderConfirmation;