import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productsAPI } from '../services/api';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import Toast from '../components/Toast';

const ProductDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { addToCart } = useCart();
  
  const [product, setProduct] = useState(null);
  const [selectedVariant, setSelectedVariant] = useState(null);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    productsAPI.getById(id)
      .then(response => {
        const productData = response.data.data;
        setProduct(productData);
        if (productData.variants && productData.variants.length > 0) {
          setSelectedVariant(productData.variants[0]);
        }
        setLoading(false);
      })
      .catch(error => {
        console.error('Failed to load product:', error);
        setLoading(false);
      });
  }, [id]);

  const handleAddToCart = async () => {
    if (!user) {
      navigate('/login');
      return;
    }

    if (!selectedVariant) return;

    try {
      setAdding(true);
      await addToCart(selectedVariant.id, 1);
      setToast({ message: 'Added to cart successfully', type: 'success' });
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Failed to add to cart';
      setToast({ message: errorMsg, type: 'error' });
    } finally {
      setAdding(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="uppercase tracking-wider text-sm">Loading...</p>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="uppercase tracking-wider text-sm">Product not found</p>
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
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-16">
        {/* Images */}
        <div>
          <div className="aspect-[3/4] bg-gray-100 border border-black mb-4">
            {product.images && product.images[0] ? (
              <img 
                src={product.images[0].image_url} 
                alt={product.name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center">
                <span className="text-sm uppercase tracking-wider text-gray-400">No Image</span>
              </div>
            )}
          </div>
          
          {/* Thumbnail Grid */}
          {product.images && product.images.length > 1 && (
            <div className="grid grid-cols-4 gap-4">
              {product.images.slice(1, 5).map((img, idx) => (
                <div key={idx} className="aspect-square bg-gray-100 border border-black">
                  <img 
                    src={img.image_url} 
                    alt={`${product.name} ${idx + 2}`}
                    className="w-full h-full object-cover"
                  />
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Product Info */}
        <div>
          <h1 className="heading-md mb-4">{product.name}</h1>
          <p className="text-2xl mb-8">{product.base_price.toFixed(2)} TL</p>

          <div className="border-t border-black pt-8 mb-8">
            <p className="text-sm leading-relaxed">{product.description}</p>
          </div>

          {/* Size Selection */}
          {product.variants && product.variants.length > 0 && (
            <div className="mb-8">
              <h3 className="text-xs uppercase tracking-wider mb-4">Select Size</h3>
              <div className="grid grid-cols-4 gap-2">
                {Array.from(new Set(product.variants.map(v => v.size))).map(size => (
                  <button
                    key={size}
                    onClick={() => {
                      const variant = product.variants.find(v => v.size === size);
                      setSelectedVariant(variant);
                    }}
                    className={`p-3 border text-sm uppercase tracking-wider transition ${
                      selectedVariant?.size === size
                        ? 'bg-black text-white border-black'
                        : 'border-black hover:bg-gray-100'
                    }`}
                  >
                    {size}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Color Selection */}
          {selectedVariant && (
            <div className="mb-8">
              <h3 className="text-xs uppercase tracking-wider mb-4">Color: {selectedVariant.color}</h3>
              <div className="flex gap-2">
                {product.variants
                  .filter(v => v.size === selectedVariant.size)
                  .map(variant => (
                    <button
                      key={variant.id}
                      onClick={() => setSelectedVariant(variant)}
                      className={`w-10 h-10 border-2 transition ${
                        selectedVariant.id === variant.id
                          ? 'border-black'
                          : 'border-gray-300'
                      }`}
                      style={{ backgroundColor: variant.color_hex }}
                      title={variant.color}
                    />
                  ))}
              </div>
            </div>
          )}

          {/* Stock Info */}
          {selectedVariant && (
            <p className="text-xs uppercase tracking-wider mb-8 text-gray-600">
              {selectedVariant.stock_quantity > 0 
                ? `${selectedVariant.stock_quantity} in stock`
                : 'Out of stock'}
            </p>
          )}

          {/* Add to Cart */}
          <button 
            onClick={handleAddToCart}
            disabled={!selectedVariant || selectedVariant.stock_quantity === 0 || adding}
            className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {adding ? 'Adding...' : 'Add to Cart'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProductDetail;