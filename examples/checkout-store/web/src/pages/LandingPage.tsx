import { useEffect, useState, type ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import { CheckCircle, Shield, Zap, Mail, ShoppingBag } from 'lucide-react';
import { fetchProduct } from '../api/client';
import type { Product } from '../types';

function formatPrice(price: number, currency: string) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(price);
}

export default function LandingPage() {
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchProduct()
      .then(setProduct)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const scrollToFeatures = () => {
    document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' });
  };

  if (loading) {
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
          <span>Failed to load store: {error || 'unknown error'}</span>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col">
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
          <div className="flex-none hidden sm:flex">
            <button className="btn btn-ghost btn-sm" onClick={scrollToFeatures}>
              Docs
            </button>
            <a
              className="btn btn-ghost btn-sm"
              href="mailto:support@example.com"
            >
              Support
            </a>
          </div>
        </div>
      </div>

      {/* Hero */}
      <div className="hero flex-grow bg-gradient-to-br from-primary/10 via-base-100 to-secondary/10 py-16">
        <div className="container mx-auto px-4">
          <div className="hero-content flex-col gap-12 lg:flex-row-reverse">
            <div className="w-full max-w-md lg:max-w-lg">
              <div className="card bg-base-100 shadow-2xl">
                <figure className="px-6 pt-6">
                  <img
                    src={product.imageUrl}
                    alt={product.name}
                    className="h-64 w-full rounded-2xl object-cover"
                  />
                </figure>
                <div className="card-body">
                  <h2 className="card-title text-2xl">{product.name}</h2>
                  <p className="text-base-content/70">{product.description}</p>
                  <div className="mt-4 flex items-end justify-between">
                    <div>
                      <p className="text-sm text-base-content/60">One-time payment</p>
                      <p className="text-4xl font-bold text-primary">
                        {formatPrice(product.price, product.currency)}
                      </p>
                    </div>
                    <button
                      className="btn btn-primary btn-lg"
                      onClick={() => navigate('/checkout')}
                    >
                      Buy now
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div className="max-w-xl">
              <div className="badge badge-primary badge-outline mb-4">New course</div>
              <h1 className="text-5xl font-extrabold leading-tight">
                Master local payment emulation with <span className="text-primary">OpenMuara</span>
              </h1>
              <p className="py-6 text-lg text-base-content/80">
                Learn how to test billing, webhooks, and provider integrations without touching
                real money. This self-paced course covers Fawry, Stripe, and more.
              </p>
              <div className="flex gap-4">
                <button
                  className="btn btn-primary btn-lg"
                  onClick={() => navigate('/checkout')}
                >
                  Get instant access
                </button>
                <button
                  className="btn btn-outline btn-lg"
                  onClick={scrollToFeatures}
                >
                  Preview syllabus
                </button>
              </div>

              <div className="mt-10 grid grid-cols-1 gap-4 sm:grid-cols-3">
                <Feature icon={<Zap className="h-5 w-5" />} text="Hands-on labs" />
                <Feature icon={<Shield className="h-5 w-5" />} text="Offline first" />
                <Feature icon={<Mail className="h-5 w-5" />} text="Email receipts" />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* What's included */}
      <div id="features" className="bg-base-200 py-16 scroll-mt-16">
        <div className="container mx-auto px-4">
          <h2 className="mb-10 text-center text-3xl font-bold">What you will learn</h2>
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            <IncludedItem title="Fawry integration" description="Create charges, handle escape pages, and verify webhooks." />
            <IncludedItem title="Stripe Checkout" description="Build sessions, confirm payments, and test card/FPX flows." />
            <IncludedItem title="Webhook testing" description="Receive provider notifications and send confirmation emails." />
            <IncludedItem title="Local ledger" description="Inspect transactions and replay events from the dashboard." />
            <IncludedItem title="Email with Mailpit" description="Test SMTP delivery without a real mail server." />
            <IncludedItem title="Provider signatures" description="Validate HMAC, SHA256, and MD5 signatures locally." />
          </div>
        </div>
      </div>

      {/* CTA */}
      <div className="bg-primary py-16 text-primary-content">
        <div className="container mx-auto px-4 text-center">
          <h2 className="mb-4 text-3xl font-bold">Ready to start?</h2>
          <p className="mx-auto mb-8 max-w-xl text-lg opacity-90">
            Buy once, access forever. Practice real provider flows in a safe, local environment.
          </p>
          <button
            className="btn btn-secondary btn-lg"
            onClick={() => navigate('/checkout')}
          >
            Buy now for {formatPrice(product.price, product.currency)}
          </button>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-base-100 py-8 text-center text-sm text-base-content/60">
        <p>© {new Date().getFullYear()} OpenMuara Store. Built for local development.</p>
      </footer>
    </div>
  );
}

function Feature({ icon, text }: { icon: ReactNode; text: string }) {
  return (
    <div className="flex items-center gap-3 rounded-lg bg-base-100/80 p-3 shadow-sm">
      <div className="text-primary">{icon}</div>
      <span className="font-medium">{text}</span>
    </div>
  );
}

function IncludedItem({ title, description }: { title: string; description: string }) {
  return (
    <div className="card bg-base-100 shadow-sm">
      <div className="card-body">
        <div className="flex items-center gap-2">
          <CheckCircle className="h-5 w-5 text-success" />
          <h3 className="card-title text-lg">{title}</h3>
        </div>
        <p className="text-base-content/70">{description}</p>
      </div>
    </div>
  );
}
