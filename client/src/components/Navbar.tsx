import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import connectionService from '../services/connectionService';
import './Navbar.css';

const Navbar = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [pendingRequestsCount, setPendingRequestsCount] = useState(0);
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();

  // Load pending connection requests count
  useEffect(() => {
    if (isAuthenticated) {
      loadRequestsCount();
      // Refresh count every 30 seconds
      const interval = setInterval(loadRequestsCount, 30000);
      return () => clearInterval(interval);
    }
  }, [isAuthenticated]);

  const loadRequestsCount = async () => {
    try {
      const count = await connectionService.getReceivedRequestsCount();
      setPendingRequestsCount(count);
    } catch (error) {
      console.error('Error loading requests count:', error);
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/');
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <div className="navbar-brand">
          <Link to="/" className="logo">
            <span className="logo-icon">💻</span>
            <span className="logo-text">Coders' Cave</span>
          </Link>
        </div>

        <div className={`navbar-menu ${isMenuOpen ? 'active' : ''}`}>
          {isAuthenticated ? (
            <>
              {/* Authenticated user navigation */}
              <Link to="/profile/me" className="navbar-link">
                My Profile
              </Link>
              <Link to="/connections" className="navbar-link">Connections</Link>
              <Link to="/messages" className="navbar-link">Messages</Link>
              <Link to="/recommendations" className="navbar-link">Discover</Link>
              <Link to="/connections/requests" className="navbar-link navbar-link-with-badge">
                Requests
                {pendingRequestsCount > 0 && (
                  <span className="navbar-badge">{pendingRequestsCount}</span>
                )}
              </Link>
              {/* Add other teammate routes when ready */}
              {/* <Link to="/recommendations" className="navbar-link">
                Discover
              </Link>
              <Link to="/chat" className="navbar-link">
                Messages
              </Link> */}
            </>
          ) : (
            <>
              {/* Public navigation */}
              <Link to="/how-it-works" className="navbar-link">How it Works</Link>
            </>
          )}
        </div>

        <div className="navbar-actions">
          {isAuthenticated ? (
            <>
              <Link to="/profile/me" className="user-profile-link">
                <div className="user-greeting">
                  <span className="user-icon">👤</span>
                  <span className="user-email">{user?.email}</span>
                </div>
              </Link>
              <button
                className="btn-secondary"
                onClick={handleLogout}
              >
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login">
                <button className="btn-secondary">Sign In</button>
              </Link>
              <Link to="/signup">
                <button className="btn-primary">Get Started</button>
              </Link>
            </>
          )}
        </div>

        <button
          className="navbar-toggle"
          onClick={() => setIsMenuOpen(!isMenuOpen)}
          aria-label="Toggle menu"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>
    </nav>
  );
};

export default Navbar;