import './Features.css';

const Features = () => {
  const features = [
    {
      icon: '🔐',
      title: 'Authentication',
      description: 'Add user sign-ups and logins, securing your data with Row Level Security.',
    },
    {
      icon: '⚡',
      title: 'Smart Matching',
      description: 'Easily match users based on preferences without deploying or scaling servers.',
    },
    {
      icon: '💬',
      title: 'Realtime Chat',
      description: 'Real-time messaging and notifications powered by WebSocket connections.',
    },
  ];

  return (
    <section className="features" id="features">
      <div className="features-container">
        <div className="features-grid">
          {features.map((feature, index) => (
            <div key={index} className="feature-card">
              <div className="feature-icon">{feature.icon}</div>
              <h3 className="feature-title">{feature.title}</h3>
              <p className="feature-description">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Features;
