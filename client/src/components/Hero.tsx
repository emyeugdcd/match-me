import { Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import './Hero.css';

const Hero = () => {
  const { isAuthenticated } = useAuth();

  return (
    <section className="hero">
      <div className="hero-container">
        <div className="hero-content">
          <h1 className="hero-title">
            Connect with Coders <br />
            <span className="gradient-text">Build Together</span>
          </h1>
          <p className="hero-subtitle">
            Coders' Cave is the ultimate platform for developers. <br />
            Find mentors, collaborators, and friends based on your tech stack and interests.
          </p>
          <div className="hero-actions">
            <Link to={isAuthenticated ? "/connections" : "/signup"}>
              <button className="hero-btn-primary">Start connecting</button>
            </Link>
          </div>

          <div className="trusted-by">
            <p>Trusted by fast-growing companies worldwide</p>
            <div className="company-logos">
              <div className="logo-placeholder">🚀</div>
              <div className="logo-placeholder">💼</div>
              <div className="logo-placeholder">🎯</div>
              <div className="logo-placeholder">⚡</div>
              <div className="logo-placeholder">🌟</div>
            </div>
          </div>
        </div>
      </div>

      <div className="hero-background">
        <div className="gradient-circle gradient-circle-1"></div>
        <div className="gradient-circle gradient-circle-2"></div>
        <div className="gradient-circle gradient-circle-3"></div>
      </div>
    </section>
  );
};

export default Hero;
