import { Routes, Route } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import CheckoutPage from './pages/CheckoutPage';
import StatusPage from './pages/StatusPage';

function App() {
  return (
    <div className="min-h-screen bg-base-100 text-base-content">
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
