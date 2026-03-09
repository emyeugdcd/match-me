import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { AuthProvider } from './hooks/useAuth';
import socketService from './services/socketService';
import './App.css';
import './styles/profiles-connections.css';

// Existing pages
import Home from './pages/Home';
import Login from './pages/Login';
import Signup from './pages/Signup';
import HowItWorksPage from './pages/HowItWorksPage';

// New pages - Profile & Connections
import EditProfilePage from './pages/EditProfilePage';
import ProfilePage from './pages/ProfilePage';
import ConnectionsPage from './pages/ConnectionsPage.tsx';
import ConnectionRequestsPage from './pages/ConnectionRequestsPage';
import ChatPage from './pages/ChatPage';
import MessagesPage from './pages/MessagesPage';
import RecommendationsPage from './pages/RecommendationsPage';

// Layout wrapper component with Navbar
import Navbar from './components/Navbar';
import authService from './services/authService';

// Protected Route wrapper component
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const isAuthenticated = authService.isAuthenticated();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

// Layout for authenticated pages (with Navbar)
const AuthenticatedLayout = ({ children }: { children: React.ReactNode }) => {
  useEffect(() => {
    socketService.connect();
    return () => {
      socketService.disconnect();
    };
  }, []);

  return (
    <>
      <Navbar />
      <div style={{ paddingTop: '70px' }}>{children}</div>
    </>
  );
};

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="app">
          <Routes>
            {/* Public routes - No navbar */}
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />

            {/* Protected routes - With navbar */}
            <Route
              path="/"
              element={<Home />}
            />
            <Route
              path="/how-it-works"
              element={<HowItWorksPage />}
            />
            <Route
              path="/profile/edit"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <EditProfilePage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile/me"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <ProfilePage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile/:userId"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <ProfilePage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/connections"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <ConnectionsPage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/connections/requests"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <ConnectionRequestsPage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />

            <Route
              path="/messages"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <MessagesPage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/chat/:userId"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <ChatPage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/recommendations"
              element={
                <ProtectedRoute>
                  <AuthenticatedLayout>
                    <RecommendationsPage />
                  </AuthenticatedLayout>
                </ProtectedRoute>
              }
            />

            {/* Catch-all redirect */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;