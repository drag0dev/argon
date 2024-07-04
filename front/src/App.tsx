import React from 'react';
import './App.css';
import HomePage from './pages/HomePage';

import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';

const About = () => {
  return (
    <div className="content">
      <h1>About Rsbuild</h1>
      <p>Learn more about Rsbuild and how it can help you build better apps.</p>
    </div>
  );
};

const App = () => {
  return (
    <Router>
      <div>
        <nav>
          <Link to="/">Home</Link> | <Link to="/about">About</Link>
        </nav>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/about" element={<About />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
