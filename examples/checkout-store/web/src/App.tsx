import { useEffect, useState } from 'react';
import { Routes, Route } from 'react-router-dom';
import AnnouncementBar from './components/AnnouncementBar';
import LandingPage from './pages/LandingPage';
import CheckoutPage from './pages/CheckoutPage';
import StatusPage from './pages/StatusPage';
import { fetchConfig } from './api/client';
import type { AppConfig } from './types';

function App() {
  const [config, setConfig] = useState<AppConfig | null>(null);

  useEffect(() => {
    fetchConfig()
      .then(setConfig)
      .catch(() => {
        // If the config endpoint is unavailable, default to demo mode so the
        // banner still nudges the user to set up their .env.
        setConfig({ demoMode: true, providers: {} as AppConfig['providers'] });
      });
  }, []);

  return (
    <div className="min-h-screen bg-base-100 text-base-content">
      <AnnouncementBar demoMode={config?.demoMode ?? false} />
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/checkout" element={<CheckoutPage />} />
        <Route path="/success" element={<StatusPage variant="success" />} />
        <Route path="/cancel" element={<StatusPage variant="cancel" />} />
      </Routes>
    </div>
  );
}

export default App;
