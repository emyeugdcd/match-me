import Navbar from '../components/Navbar';
import './HowItWorksPage.css';

const HowItWorksPage = () => {
    return (
        <>
            <Navbar />
            <div className="how-it-works-page">
                <div className="container">
                    <div className="content-wrapper">
                        <h1 className="page-title">How Coders' Cave Works</h1>

                        <div className="intro-section">
                            <p className="lead-text">
                                Coders' Cave is designed with one simple mission: <strong>to connect developers through code.</strong>
                            </p>
                            <p>
                                Whether you're looking for a mentor, a hackathon partner, or just a friend to discuss the latest tech trends with,
                                we make it easy to find people who speak your language—literally.
                            </p>
                        </div>

                        <div className="steps-grid">
                            <div className="step-card">
                                <div className="step-number">1</div>
                                <h3>Create Your Profile</h3>
                                <p>Sign up and tell us about yourself. Most importantly, list your <strong>Programming Languages</strong>. Let others know if you're into Python, Go, Rust, or JavaScript.</p>
                            </div>

                            <div className="step-card">
                                <div className="step-number">2</div>
                                <h3>Discover Peers</h3>
                                <p>Our recommendation engine suggests other developers based on shared tech stacks, interests, and location. No more random connections—find people who actually get your code jokes.</p>
                            </div>

                            <div className="step-card">
                                <div className="step-number">3</div>
                                <h3>Connect & Collaborate</h3>
                                <p>Send connection requests and start chatting. Build projects together, share knowledge, or just expand your professional network in a meaningful way.</p>
                            </div>
                        </div>

                        <div className="cta-section">
                            <h2>Ready to find your squad?</h2>
                            <a href="/signup" className="cta-button">Join Coders' Cave Today</a>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};

export default HowItWorksPage;
