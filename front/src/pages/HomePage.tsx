import React from 'react';
import HomeFeed from '../components/HomeFeed';

const HomePage = () => {
  return (
    <section className="section">
      <div className="container">
        <div style={{ minHeight: '100vh' }}>
          <section className="hero is-secondary">
            <div className="hero-body">
              <div className="container">
                <h1 className="title is-1">
                  Argon
                </h1>
                <h2 className="subtitle is-3">
                  Stable. Noble. Entertainment.
                </h2>
              </div>
            </div>
          </section>
          <HomeFeed />
        </div>
      </div>
    </section>
  );
};

export default HomePage;
