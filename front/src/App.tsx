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
import DocumentationPage from './pages/DocumentationPage';
import RegisterPage from './pages/RegisterPage';
import LoginPage from './pages/LoginPage';

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Amplify } from 'aws-amplify';

const APP_POOL_ID = process.env.APP_POOL_ID;
const POOL_CLIENT_ID = process.env.POOL_CLIENT_ID;

const App = () => {
      Amplify.configure({
          aws_project_region: 'eu-central-1',
          aws_user_pools_id: APP_POOL_ID,
          aws_user_pools_web_client_id: POOL_CLIENT_ID,
      });
  return (
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path="/" element={<HomePage />} />

          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route path="/about" element={<AboutPage />} />
          <Route path="/docs" element={<DocumentationPage />} />
          <Route path="/movie/add" element={<AddMoviePage />} />
          <Route path="/tvshow/add" element={<AddTVShowPage />} />
          <Route path="/movie/delete" element={<EditOrDeleteMoviePage />} />
          <Route path="/movie/edit" element={<EditOrDeleteMoviePage />} />
          <Route path="/tvshow/delete" element={<EditOrDeleteTVShowPage />} />
          <Route path="/tvshow/edit" element={<EditOrDeleteTVShowPage />} />
          <Route path="/movie/:uuid/details" element={<WatchPage />} />
          <Route path="/tvshow/:uuid/details" element={<TVShowDetails />} />
          <Route path="/tvshow/:uuid/watch/:seasonId/:episodeId" element={<WatchPage />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
