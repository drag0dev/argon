import React from 'react';
import './App.css';
import HomePage from './pages/HomePage';
import AboutPage from './pages/AboutPage';
import NavBar from './components/NavBar';
import AddMoviePage from './pages/AddMoviePage';
import AddTVShowPage from './pages/AddTVShowPage';
import EditOrDeleteMoviePage from './pages/EditOrDeleteMoviePage';

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

const App = () => {
  return (
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/about" element={<AboutPage />} />
          <Route path="/movie/add" element={<AddMoviePage />} />
          <Route path="/tvshow/add" element={<AddTVShowPage />} />
          <Route path="/movie/delete" element={<EditOrDeleteMoviePage />} />
          <Route path="/movie/edit" element={<EditOrDeleteMoviePage />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
