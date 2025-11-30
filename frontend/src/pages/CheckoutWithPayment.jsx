import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import { useCart } from '../context/CartContext';
import { addressAPI, orderAPI } from '../services/api';
import PaymentForm from '../components/PaymentForm';
import Toast from '../components/Toast';

// Load Stripe with your publishable key
const stripePromise = loadStripe('pk_test_51SZGwqF3TJGbZvctotzYkCiPIlThaylELuk8bJ9EpphaKDJsGwfc9gHnx9vT0cjDuNjOV3otz3RIZ9nWlvru4In400d4qxKkxd');

const CheckoutWithPayment = () => {
  const navigate = useNavigate();
  const { cart, fetchCart } = useCart();
  const [step, setStep] = useState(1); // 1: Address, 2: Payment
  const [addresses, setAddresses] = useState([]);
  const [selectedAddress, setSelectedAddress] = useState(null);
  const [orderId, setOrderId] = useState(null);
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    console.log('Cart data:', cart); // Debug
    if (!cart || cart.items?.length === 0) {
      navigate('/cart');
    }
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

  const handleCreateOrder = async () => {
    if (!selectedAddress) {
      setToast({ type: 'error', message: 'Please select a shipping address' });
      return;
    }

    setLoading(true);
    try {
      const response = await orderAPI.createOrder({
        address_id: selectedAddress,
        payment_method: 'stripe',
      });

      setOrderId(response.data.data.id);
      setStep(2); // Move to payment step
    } catch (error) {
      setToast({
        type: 'error',
        message: error.response?.data?.error || 'Failed to create order',
      });
    } finally {
      setLoading(false);
    }
  };

  const handlePaymentSuccess = async (paymentIntent) => {
    setToast({ type: 'success', message: 'Payment successful!' });
    await fetchCart(); // Clear cart
    setTimeout(() => {
      navigate(`/order-confirmation/${orderId}`);
    }, 1500);
  };

  const handlePaymentError = (error) => {
    setToast({ type: 'error', message: error.message || 'Payment failed' });
  };

  if (!cart || cart.items?.length === 0) {
    return null;
  }

  return (
    <div className="min-h-screen bg-white py-12">
      <div className="container mx-auto px-8 max-w-4xl">
        {/* Progress Steps */}
        <div className="mb-12">
          <div className="flex items-center justify-center space-x-4">
            <div className={`flex items-center ${step >= 1 ? 'opacity-100' : 'opacity-50'}`}>
              <div className="w-8 h-8 border-2 border-black flex items-center justify-center">
                {step > 1 ? '✓' : '1'}
              </div>
              <span className="ml-2 uppercase tracking-wider text-sm">ADDRESS</span>
            </div>
            <div className="w-16 border-t-2 border-black"></div>
            <div className={`flex items-center ${step >= 2 ? 'opacity-100' : 'opacity-50'}`}>
              <div className="w-8 h-8 border-2 border-black flex items-center justify-center">
                2
              </div>
              <span className="ml-2 uppercase tracking-wider text-sm">PAYMENT</span>
            </div>
          </div>
        </div>

        {/* Step 1: Address Selection */}
        {step === 1 && (
          <div className="space-y-8">
            <h1 className="text-2xl uppercase tracking-widest mb-8">SELECT SHIPPING ADDRESS</h1>

            {addresses.length === 0 ? (
              <div className="border border-black p-8 text-center">
                <p className="mb-4">NO ADDRESSES FOUND</p>
                <button
                  onClick={() => navigate('/profile')}
                  className="bg-black text-white px-6 py-2 uppercase tracking-wider"
                >
                  ADD ADDRESS
                </button>
              </div>
            ) : (
              <div className="space-y-4">
                {addresses.map((address) => (
                  <div
                    key={address.id}
                    onClick={() => setSelectedAddress(address.id)}
                    className={`border-2 p-6 cursor-pointer hover:bg-gray-50 ${
                      selectedAddress === address.id ? 'border-black bg-gray-50' : 'border-gray-300'
                    }`}
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-bold uppercase tracking-wider mb-2">{address.title}</p>
                        <p className="text-sm">{address.full_name}</p>
                        <p className="text-sm">{address.address_line1}</p>
                        {address.address_line2 && <p className="text-sm">{address.address_line2}</p>}
                        <p className="text-sm">
                          {address.city}, {address.state} {address.postal_code}
                        </p>
                        <p className="text-sm">{address.country}</p>
                        <p className="text-sm mt-2">{address.phone}</p>
                      </div>
                      <div>
                        {selectedAddress === address.id && (
                          <div className="w-6 h-6 bg-black text-white flex items-center justify-center">
                            ✓
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* Order Summary */}
            <div className="border border-black p-6 mt-8">
              <h2 className="text-lg uppercase tracking-wider mb-4">ORDER SUMMARY</h2>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>ITEMS:</span>
                  <span>{cart?.total_items || 0}</span>
                </div>
                <div className="flex justify-between font-bold text-base pt-4 border-t border-black">
                  <span>TOTAL:</span>
                  <span>${(cart?.total_price || 0).toFixed(2)}</span>
                </div>
              </div>
            </div>

            {/* Continue Button */}
            <button
              onClick={handleCreateOrder}
              disabled={!selectedAddress || loading}
              className="w-full bg-black text-white py-4 uppercase tracking-wider hover:opacity-80 disabled:opacity-50"
            >
              {loading ? 'CREATING ORDER...' : 'CONTINUE TO PAYMENT'}
            </button>
          </div>
        )}

        {/* Step 2: Payment */}
        {step === 2 && orderId && cart?.total_price && (
          <div className="space-y-8">
            <h1 className="text-2xl uppercase tracking-widest mb-8">PAYMENT</h1>

            <Elements stripe={stripePromise}>
              <PaymentForm
                amount={cart.total_price || 0}
                orderId={orderId}
                onSuccess={handlePaymentSuccess}
                onError={handlePaymentError}
              />
            </Elements>

            {/* Back Button */}
            <button
              onClick={() => setStep(1)}
              className="w-full border border-black py-3 uppercase tracking-wider hover:bg-gray-100"
            >
              BACK TO ADDRESS
            </button>
          </div>
        )}

        {/* Toast Notifications */}
        {toast && (
          <Toast
            type={toast.type}
            message={toast.message}
            onClose={() => setToast(null)}
          />
        )}
      </div>
    </div>
  );
};

export default CheckoutWithPayment;