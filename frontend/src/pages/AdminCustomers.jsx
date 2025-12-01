import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || '185.26.144.214:4501';

const AdminCustomers = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Redirect if not admin
    if (user && user.role !== 'admin') {
      navigate('/');
      return;
    }

    // Fetch customers
    const token = localStorage.getItem('token');
    axios.get(`${API_URL}/admin/customers`, {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then(response => {
        setCustomers(response.data.data || []);
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load customers:', error);
        setLoading(false);
      });
  }, [user, navigate]);

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
        <h1 className="heading-lg">Customer Management</h1>
        <button
          onClick={() => navigate('/admin')}
          className="btn-secondary"
        >
          Back to Dashboard
        </button>
      </div>

      <div className="mb-6 text-sm uppercase tracking-wider text-gray-600">
        Total Customers: {customers.length}
      </div>

      {/* Customers Table */}
      <div className="border border-black">
        <table className="w-full">
          <thead className="border-b border-black bg-gray-50">
            <tr>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">ID</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Name</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Email</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Phone</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Role</th>
              <th className="px-6 py-4 text-left text-xs uppercase tracking-wider">Joined</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {customers.length === 0 ? (
              <tr>
                <td colSpan="6" className="px-6 py-8 text-center text-gray-500">
                  No customers found
                </td>
              </tr>
            ) : (
              customers.map((customer) => (
                <tr key={customer.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm">{customer.id}</td>
                  <td className="px-6 py-4 text-sm">
                    {customer.first_name} {customer.last_name}
                  </td>
                  <td className="px-6 py-4 text-sm">{customer.email}</td>
                  <td className="px-6 py-4 text-sm">{customer.phone || '-'}</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs uppercase tracking-wider px-2 py-1 ${
                      customer.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-gray-100 text-gray-800'
                    }`}>
                      {customer.role}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(customer.created_at).toLocaleDateString()}
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

export default AdminCustomers;