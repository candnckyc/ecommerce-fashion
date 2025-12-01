import React, { useState } from 'react';
import { CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || '185.26.144.214:4501';

const PaymentForm = ({ amount = 0, orderId, onSuccess, onError }) => {
  const stripe = useStripe();
  const elements = useElements();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Safety check
  if (!amount || amount <= 0) {
    return (
      <div className="border border-red-500 p-6 text-center">
        <p className="text-red-500 mb-4">Invalid payment amount</p>
        <p className="text-sm">Please go back and try again.</p>
      </div>
    );
  }

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!stripe || !elements) {
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Step 1: Create payment intent
      const token = localStorage.getItem('token');
      const { data } = await axios.post(
        `${API_URL}/payment/create-intent`,
        {
          amount: Math.round(amount * 100), // Convert to cents
          currency: 'usd',
          order_id: orderId,
        },
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const { client_secret, payment_intent_id } = data.data;

      // Step 2: Confirm payment with card
      const result = await stripe.confirmCardPayment(client_secret, {
        payment_method: {
          card: elements.getElement(CardElement),
        },
      });

      if (result.error) {
        setError(result.error.message);
        if (onError) onError(result.error);
      } else {
        // Step 3: Payment successful - confirm on backend
        await axios.post(
          `${API_URL}/payment/confirm`,
          {
            payment_intent_id,
            order_id: orderId,
          },
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );

        if (onSuccess) onSuccess(result.paymentIntent);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Payment failed');
      if (onError) onError(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Card Element */}
      <div className="border border-black p-4">
        <CardElement
          options={{
            style: {
              base: {
                fontSize: '16px',
                color: '#000',
                '::placeholder': {
                  color: '#999',
                },
              },
              invalid: {
                color: '#ff0000',
              },
            },
          }}
        />
      </div>

      {/* Error Display */}
      {error && (
        <div className="bg-red-50 border border-red-500 px-4 py-3 text-sm">
          {error}
        </div>
      )}

      {/* Amount Display */}
      <div className="text-sm uppercase tracking-wider">
        Total Amount: <span className="font-bold">${amount.toFixed(2)}</span>
      </div>

      {/* Submit Button */}
      <button
        type="submit"
        disabled={!stripe || loading}
        className="w-full bg-black text-white py-3 uppercase tracking-wider hover:opacity-80 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? 'PROCESSING...' : `PAY $${amount.toFixed(2)}`}
      </button>

      {/* Test Card Info */}
      <div className="text-xs text-gray-500 mt-4 border-t border-gray-200 pt-4">
        <p className="font-bold mb-2">TEST CARD:</p>
        <p>Number: 4242 4242 4242 4242</p>
        <p>Expiry: Any future date (e.g., 12/25)</p>
        <p>CVC: Any 3 digits</p>
      </div>
    </form>
  );
};

export default PaymentForm;