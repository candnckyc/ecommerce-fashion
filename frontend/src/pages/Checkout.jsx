import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { addressAPI, orderAPI } from '../services/api';
import Toast from '../components/Toast';

const Checkout = () => {
  const navigate = useNavigate();
  const { cart, fetchCart } = useCart();
  const [addresses, setAddresses] = useState([]);
  const [selectedAddress, setSelectedAddress] = useState(null);
  const [paymentMethod, setPaymentMethod] = useState('credit_card');
  const [showAddressForm, setShowAddressForm] = useState(false);
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  const [newAddress, setNewAddress] = useState({
    title: '',
    full_name: '',
    phone: '',
    address_line1: '',
    address_line2: '',
    city: '',
    state: '',
    postal_code: '',
    country: 'Turkey',
    is_default: false,
  });

  useEffect(() => {
    fetchAddresses();
  }, []);

  const fetchAddresses = async () => {
    try {
      const response = await addressAPI.getAddresses();
      const addrs = response.data.data || [];
      setAddresses(addrs);
      if (addrs.length > 0) {
        setSelectedAddress(addrs[0].id);
      }
    } catch (error) {
      console.error('Failed to load addresses:', error);
    }
  };

  const handleAddAddress = async (e) => {
    e.preventDefault();
    try {
      await addressAPI.createAddress(newAddress);
      setShowAddressForm(false);
      fetchAddresses();
      setToast({ message: 'Address added', type: 'success' });
    } catch (error) {
      setToast({ message: 'Failed to add address', type: 'error' });
    }
  };

  const handlePlaceOrder = async () => {
    if (!selectedAddress) {
      setToast({ message: 'Please select an address', type: 'error' });
      return;
    }

    try {
      setLoading(true);
      const response = await orderAPI.createOrder({
        address_id: selectedAddress,
        payment_method: paymentMethod,
        notes: '',
      });
      const order = response.data.data;
      
      // Refresh cart (backend already cleared it)
      await fetchCart();
      
      navigate(`/order-confirmation/${order.id}`);
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Failed to place order';
      setToast({ message: errorMsg, type: 'error' });
    } finally {
      setLoading(false);
    }
  };

  if (cart.items?.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="heading-md mb-8">Your cart is empty</h1>
          <button onClick={() => navigate('/products')} className="btn-primary">
            Continue Shopping
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-8 py-12">
      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}

      <h1 className="heading-lg mb-12">Checkout</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
        {/* Left Column - Address & Payment */}
        <div className="md:col-span-2 space-y-12">
          {/* Shipping Address */}
          <div>
            <h2 className="text-sm uppercase tracking-wider mb-6">Shipping Address</h2>

            {addresses.length === 0 && !showAddressForm && (
              <button
                onClick={() => setShowAddressForm(true)}
                className="btn-secondary"
              >
                Add Address
              </button>
            )}

            {addresses.length > 0 && !showAddressForm && (
              <div className="space-y-4">
                {addresses.map((addr) => (
                  <label
                    key={addr.id}
                    className="block border border-black p-6 cursor-pointer hover:border-2"
                  >
                    <input
                      type="radio"
                      name="address"
                      value={addr.id}
                      checked={selectedAddress === addr.id}
                      onChange={() => setSelectedAddress(addr.id)}
                      className="mr-4"
                    />
                    <div className="inline-block text-sm">
                      <p className="font-bold">{addr.title} - {addr.full_name}</p>
                      <p>{addr.address_line1}</p>
                      {addr.address_line2 && <p>{addr.address_line2}</p>}
                      <p>{addr.city}, {addr.state} {addr.postal_code}</p>
                      <p>{addr.country}</p>
                      <p className="mt-2">{addr.phone}</p>
                    </div>
                  </label>
                ))}
                <button
                  onClick={() => setShowAddressForm(true)}
                  className="btn-secondary mt-4"
                >
                  Add New Address
                </button>
              </div>
            )}

            {showAddressForm && (
              <form onSubmit={handleAddAddress} className="border border-black p-6 space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <input
                    type="text"
                    placeholder="Title (Home, Work)"
                    className="input-field"
                    value={newAddress.title}
                    onChange={(e) => setNewAddress({ ...newAddress, title: e.target.value })}
                    required
                  />
                  <input
                    type="text"
                    placeholder="Full Name"
                    className="input-field"
                    value={newAddress.full_name}
                    onChange={(e) => setNewAddress({ ...newAddress, full_name: e.target.value })}
                    required
                  />
                </div>
                <input
                  type="tel"
                  placeholder="Phone"
                  className="input-field"
                  value={newAddress.phone}
                  onChange={(e) => setNewAddress({ ...newAddress, phone: e.target.value })}
                  required
                />
                <input
                  type="text"
                  placeholder="Address Line 1"
                  className="input-field"
                  value={newAddress.address_line1}
                  onChange={(e) => setNewAddress({ ...newAddress, address_line1: e.target.value })}
                  required
                />
                <input
                  type="text"
                  placeholder="Address Line 2 (Optional)"
                  className="input-field"
                  value={newAddress.address_line2}
                  onChange={(e) => setNewAddress({ ...newAddress, address_line2: e.target.value })}
                />
                <div className="grid grid-cols-3 gap-4">
                  <input
                    type="text"
                    placeholder="City"
                    className="input-field"
                    value={newAddress.city}
                    onChange={(e) => setNewAddress({ ...newAddress, city: e.target.value })}
                    required
                  />
                  <input
                    type="text"
                    placeholder="State"
                    className="input-field"
                    value={newAddress.state}
                    onChange={(e) => setNewAddress({ ...newAddress, state: e.target.value })}
                  />
                  <input
                    type="text"
                    placeholder="Postal Code"
                    className="input-field"
                    value={newAddress.postal_code}
                    onChange={(e) => setNewAddress({ ...newAddress, postal_code: e.target.value })}
                  />
                </div>
                <div className="flex gap-4">
                  <button type="submit" className="btn-primary">Save Address</button>
                  <button
                    type="button"
                    onClick={() => setShowAddressForm(false)}
                    className="btn-secondary"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            )}
          </div>

          {/* Payment Method */}
          <div>
            <h2 className="text-sm uppercase tracking-wider mb-6">Payment Method</h2>
            <div className="space-y-4">
              <label className="block border border-black p-6 cursor-pointer hover:border-2">
                <input
                  type="radio"
                  name="payment"
                  value="credit_card"
                  checked={paymentMethod === 'credit_card'}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  className="mr-4"
                />
                <span className="text-sm uppercase tracking-wider">Credit Card</span>
              </label>
              <label className="block border border-black p-6 cursor-pointer hover:border-2">
                <input
                  type="radio"
                  name="payment"
                  value="cash_on_delivery"
                  checked={paymentMethod === 'cash_on_delivery'}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  className="mr-4"
                />
                <span className="text-sm uppercase tracking-wider">Cash on Delivery</span>
              </label>
            </div>
          </div>
        </div>

        {/* Right Column - Order Summary */}
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

            <button
              onClick={handlePlaceOrder}
              disabled={loading || !selectedAddress}
              className="btn-primary w-full disabled:opacity-50"
            >
              {loading ? 'Placing Order...' : 'Place Order'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Checkout;