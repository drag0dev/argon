import React from 'react';
import HomeFeed from '../components/HomeFeed';

const HomePage = () => {
  return (
    <div style={{ minHeight: '100vh' }}>
      <section className="hero is-primary is-fullheight">
        <div className="hero-body">
          <div className="container">
            <h1 className="title is-1">Welcome to the Movie/TV Show Tracker</h1>
            <h2 className="subtitle is-3">Keep track of your favorite movies and TV shows</h2>
            <HomeFeed />
          </div>
        </div>
      </section>
    </div>
  );
};

export default HomePage;
