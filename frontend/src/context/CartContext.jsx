import React, { createContext, useState, useContext, useEffect } from 'react';
import { cartAPI } from '../services/api';
import { useAuth } from './AuthContext';

const CartContext = createContext();

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within CartProvider');
  }
  return context;
};

export const CartProvider = ({ children }) => {
  const [cart, setCart] = useState({ items: [], total_items: 0, total_price: 0 });
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();

  // Fetch cart when user logs in
  useEffect(() => {
    if (user) {
      fetchCart();
    } else {
      setCart({ items: [], total_items: 0, total_price: 0 });
    }
  }, [user]);

  const fetchCart = async () => {
    try {
      setLoading(true);
      const response = await cartAPI.getCart();
      setCart(response.data.data);
    } catch (error) {
      console.error('Failed to fetch cart:', error);
    } finally {
      setLoading(false);
    }
  };

  const addToCart = async (productVariantId, quantity) => {
    try {
      const response = await cartAPI.addToCart({ product_variant_id: productVariantId, quantity });
      setCart(response.data.data);
      return response.data.data;
    } catch (error) {
      throw error;
    }
  };

  const updateCartItem = async (cartItemId, quantity) => {
    try {
      const response = await cartAPI.updateCartItem(cartItemId, { quantity });
      setCart(response.data.data);
    } catch (error) {
      throw error;
    }
  };

  const removeFromCart = async (cartItemId) => {
    try {
      const response = await cartAPI.removeFromCart(cartItemId);
      setCart(response.data.data);
    } catch (error) {
      throw error;
    }
  };

  return (
    <CartContext.Provider value={{ 
      cart, 
      loading, 
      addToCart, 
      updateCartItem, 
      removeFromCart, 
      fetchCart 
    }}>
      {children}
    </CartContext.Provider>
  );
};
