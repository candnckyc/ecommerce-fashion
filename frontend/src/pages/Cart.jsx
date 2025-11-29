import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import Toast from '../components/Toast';

const Cart = () => {
  const { cart, loading, updateCartItem, removeFromCart } = useCart();
  const [toast, setToast] = useState(null);

  const handleQuantityChange = async (itemId, newQuantity, maxStock) => {
    if (newQuantity < 1) return;
    
    // Check against available stock
    if (newQuantity > maxStock) {
      setToast({ message: `Only ${maxStock} items available in stock`, type: 'error' });
      return;
    }
    
    try {
      await updateCartItem(itemId, newQuantity);
      setToast({ message: 'Cart updated', type: 'success' });
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Failed to update quantity';
      setToast({ message: errorMsg, type: 'error' });
    }
  };

  const handleRemove = async (itemId) => {
    try {
      await removeFromCart(itemId);
      setToast({ message: 'Item removed', type: 'success' });
    } catch (error) {
      setToast({ message: 'Failed to remove item', type: 'error' });
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="uppercase tracking-wider text-sm">Loading...</p>
      </div>
    );
  }

  if (!cart.items || cart.items.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="heading-md mb-8">Your cart is empty</h1>
          <Link to="/products" className="btn-primary">
            Continue Shopping
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-8 py-12">
      {toast && (
        <Toast 
          message={toast.message} 
          type={toast.type} 
          onClose={() => setToast(null)} 
        />
      )}
      
      <h1 className="heading-lg mb-12">Shopping Cart</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
        {/* Cart Items */}
        <div className="md:col-span-2">
          <div className="space-y-px bg-black">
            {cart.items.map(item => (
              <div key={item.id} className="bg-white p-6">
                <div className="flex gap-6">
                  {/* Image */}
                  <div className="w-32 h-40 bg-gray-100 border border-black flex-shrink-0">
                    {item.product?.images?.[0] ? (
                      <img 
                        src={item.product.images[0].image_url}
                        alt={item.product.name}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center">
                        <span className="text-xs uppercase">No Image</span>
                      </div>
                    )}
                  </div>

                  {/* Details */}
                  <div className="flex-1">
                    <Link 
                      to={`/products/${item.product.id}`}
                      className="text-sm uppercase tracking-wider hover:underline"
                    >
                      {item.product?.name}
                    </Link>
                    <p className="text-xs mt-2 text-gray-600">
                      Size: {item.variant?.size} | Color: {item.variant?.color}
                    </p>
                    <p className="text-sm mt-2">{item.variant?.final_price?.toFixed(2)} TL</p>

                    {/* Quantity Controls */}
                    <div className="flex items-center gap-2 mt-4">
                      <button
                        onClick={() => handleQuantityChange(item.id, item.quantity - 1, item.variant?.stock_quantity)}
                        className="w-8 h-8 border border-black hover:bg-black hover:text-white transition"
                      >
                        −
                      </button>
                      <span className="w-12 text-center text-sm">{item.quantity}</span>
                      <button
                        onClick={() => handleQuantityChange(item.id, item.quantity + 1, item.variant?.stock_quantity)}
                        className="w-8 h-8 border border-black hover:bg-black hover:text-white transition"
                      >
                        +
                      </button>
                    </div>
                  </div>

                  {/* Remove Button */}
                  <button
                    onClick={() => handleRemove(item.id)}
                    className="text-sm uppercase tracking-wider hover:opacity-50 transition self-start"
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Order Summary */}
        <div className="md:col-span-1">
          <div className="border border-black p-6 sticky top-8">
            <h2 className="text-sm uppercase tracking-wider mb-6">Order Summary</h2>

            <div className="space-y-4 mb-6">
              <div className="flex justify-between text-sm">
                <span>Subtotal ({cart.total_items} items)</span>
                <span>{cart.total_price?.toFixed(2)} TL</span>
              </div>
              <div className="flex justify-between text-sm">
                <span>Shipping</span>
                <span>Free</span>
              </div>
            </div>

            <div className="border-t border-black pt-4 mb-6">
              <div className="flex justify-between text-lg">
                <span className="uppercase tracking-wider">Total</span>
                <span>{cart.total_price?.toFixed(2)} TL</span>
              </div>
            </div>

            <button className="btn-primary w-full mb-4">
              Proceed to Checkout
            </button>

            <Link 
              to="/products"
              className="block text-center text-sm uppercase tracking-wider hover:underline"
            >
              Continue Shopping
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Cart;