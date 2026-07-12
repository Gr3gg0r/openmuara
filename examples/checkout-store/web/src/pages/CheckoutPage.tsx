import { useEffect, useState, type ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, CreditCard, Landmark, Loader2, ShoppingBag } from 'lucide-react';
import { createCheckout, fetchProduct } from '../api/client';
import type { Product } from '../types';

function formatPrice(price: number, currency: string) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(price);
}

export default function CheckoutPage() {
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [method, setMethod] = useState<'fawry' | 'stripe'>('fawry');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [productLoading, setProductLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchProduct()
      .then(setProduct)
      .catch((err) => setError(err.message))
      .finally(() => setProductLoading(false));
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const res = await createCheckout({ method, name, email });
      window.location.href = res.redirectUrl;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Checkout failed');
      setLoading(false);
    }
  };

  if (productLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="alert alert-error max-w-md">
          <span>Failed to load checkout: {error || 'unknown error'}</span>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Navbar */}
      <div className="navbar bg-base-100 shadow-sm">
        <div className="container mx-auto px-4">
          <div className="flex-1">
            <button
              className="btn btn-ghost text-xl font-bold"
              onClick={() => navigate('/')}
            >
              <ShoppingBag className="h-6 w-6 text-primary" />
              OpenMuara Store
            </button>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-12">
        <button className="btn btn-ghost mb-6" onClick={() => navigate('/')}>
          <ArrowLeft className="h-4 w-4" />
          Back to store
        </button>

        <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
          {/* Order summary */}
          <div>
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body">
                <h2 className="card-title text-2xl">Order summary</h2>
                <div className="mt-4 flex gap-4">
                  <img
                    src={product.imageUrl}
                    alt={product.name}
                    className="h-24 w-24 rounded-lg object-cover"
                  />
                  <div>
                    <h3 className="font-bold">{product.name}</h3>
                    <p className="text-sm text-base-content/70">{product.description}</p>
                  </div>
                </div>
                <div className="divider"></div>
                <div className="flex justify-between text-lg">
                  <span>Subtotal</span>
                  <span>{formatPrice(product.price, product.currency)}</span>
                </div>
                <div className="flex justify-between text-base-content/70">
                  <span>Tax</span>
                  <span>Included</span>
                </div>
                <div className="divider"></div>
                <div className="flex justify-between text-2xl font-bold">
                  <span>Total</span>
                  <span className="text-primary">
                    {formatPrice(product.price, product.currency)}
                  </span>
                </div>
              </div>
            </div>

            <div className="card mt-6 bg-base-100 shadow-sm">
              <div className="card-body">
                <h3 className="card-title text-lg">Secure checkout</h3>
                <p className="text-sm text-base-content/70">
                  Payments are processed through OpenMuara's local Fawry and Stripe emulators.
                  No real money is charged.
                </p>
              </div>
            </div>
          </div>

          {/* Payment form */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body">
              <h2 className="card-title text-2xl">Payment details</h2>
              <form onSubmit={handleSubmit} className="mt-4 space-y-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Full name</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered w-full"
                    placeholder="Ahmad Ali"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Email</span>
                  </label>
                  <input
                    type="email"
                    className="input input-bordered w-full"
                    placeholder="ahmad@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt">Your receipt will be sent here</span>
                  </label>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Payment method</span>
                  </label>
                  <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                    <MethodCard
                      selected={method === 'fawry'}
                      onClick={() => setMethod('fawry')}
                      icon={<Landmark className="h-6 w-6" />}
                      title="Fawry"
                      description="Pay via OpenMuara Fawry emulator"
                    />
                    <MethodCard
                      selected={method === 'stripe'}
                      onClick={() => setMethod('stripe')}
                      icon={<CreditCard className="h-6 w-6" />}
                      title="Stripe"
                      description="Card via OpenMuara Stripe emulator"
                    />
                  </div>
                </div>

                {error && (
                  <div className="alert alert-error">
                    <span>{error}</span>
                  </div>
                )}

                <button
                  type="submit"
                  className="btn btn-primary btn-block btn-lg"
                  disabled={loading}
                >
                  {loading ? (
                    <>
                      <Loader2 className="h-5 w-5 animate-spin" />
                      Processing...
                    </>
                  ) : (
                    `Pay ${formatPrice(product.price, product.currency)}`
                  )}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function MethodCard({
  selected,
  onClick,
  icon,
  title,
  description,
}: {
  selected: boolean;
  onClick: () => void;
  icon: ReactNode;
  title: string;
  description: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`flex items-start gap-3 rounded-xl border p-4 text-left transition-all ${
        selected
          ? 'border-primary bg-primary/10 ring-1 ring-primary'
          : 'border-base-300 bg-base-100 hover:border-primary/50'
      }`}
    >
      <div className={selected ? 'text-primary' : 'text-base-content/60'}>{icon}</div>
      <div>
        <div className="font-bold">{title}</div>
        <div className="text-sm text-base-content/70">{description}</div>
      </div>
    </button>
  );
}
