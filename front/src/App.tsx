import React from 'react';
import './App.css';
import HomePage from './pages/HomePage';
import AboutPage from './pages/AboutPage';
import NavBar from './components/NavBar';
import AddMoviePage from './pages/AddMoviePage';
import AddTVShowPage from './pages/AddTVShowPage';
import EditOrDeleteMoviePage from './pages/EditOrDeleteMoviePage';
import EditOrDeleteTVShowPage from './pages/EditOrDeleteTVShowPage';
import WatchPage from './pages/WatchPage';
import TVShowDetails from './components/TVShowDetails';
import RegisterPage from './pages/RegisterPage';
import LoginPage from './pages/LoginPage';

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

const App = () => {
  return (
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path="/" element={<HomePage />} />

          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route path="/about" element={<AboutPage />} />
          <Route path="/movie/add" element={<AddMoviePage />} />
          <Route path="/tvshow/add" element={<AddTVShowPage />} />
          <Route path="/movie/delete" element={<EditOrDeleteMoviePage />} />
          <Route path="/movie/edit" element={<EditOrDeleteMoviePage />} />
          <Route path="/tvshow/delete" element={<EditOrDeleteTVShowPage />} />
          <Route path="/tvshow/edit" element={<EditOrDeleteTVShowPage />} />
          <Route path="/movie/:id/details" element={<WatchPage />} />
          <Route path="/tvshow/:id/details" element={<TVShowDetails />} />
          <Route path="/tvshow/:id/watch/:seasonId/:episodeId" element={<WatchPage />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
