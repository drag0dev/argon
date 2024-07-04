import React from 'react';
import AddMovieForm from '../components/AddMovieForm'

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUser } from '@fortawesome/free-solid-svg-icons';

const AddMoviePage = () => {
  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Add New Movie</h1>
        <AddMovieForm />
      </div>
    </section>
  );
}

export default AddMoviePage;
