import React, { useState, useEffect } from 'react';
import { API_URL } from '../../config.ts';
import EditMovieForm from './EditMovieForm';

const API_URL = process.env.API_URL;

const MovieTable = () => {
  const [movies, setMovies] = useState([]);
  const [editingMovie, setEditingMovie] = useState(null);

  const dummyMovies = [
    {
      id: 1,
      title: 'Inception',
      description:
        'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.',
      genres: ['Action', 'Adventure', 'Sci-Fi'],
      actors: ['Leonardo DiCaprio', 'Joseph Gordon-Levitt', 'Ellen Page'],
      directors: ['Christopher Nolan'],
    },
    {
      id: 2,
      title: 'Interstellar',
      description:
        "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
      genres: ['Adventure', 'Drama', 'Sci-Fi'],
      actors: ['Matthew McConaughey', 'Anne Hathaway', 'Jessica Chastain'],
      directors: ['Christopher Nolan'],
    },
    {
      id: 3,
      title: 'The Dark Knight',
      description:
        'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.',
      genres: ['Action', 'Crime', 'Drama'],
      actors: ['Christian Bale', 'Heath Ledger', 'Aaron Eckhart'],
      directors: ['Christopher Nolan'],
    },
  ];

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch(`${API_URL}/api/movies`);
        const data = await response.json();
        setMovies(data);
      } catch (error) {
        console.error('Failed to fetch movies:', error);
      }
    };

    setMovies(dummyMovies);

    // TODO: actually integrate
    // fetchMovies();
  }, []);

  const handleDelete = async (movieId) => {
    const confirmed = window.confirm(
      'Are you sure you want to delete this movie?',
    );
    if (confirmed) {
      setIsLoading(true);
      try {
        await fetch(`${API_URL}/api/movies/${movieId}`, {
          method: 'DELETE',
        });
        // Remove movie from state after successful deletion
        setMovies(movies.filter((movie) => movie.id !== movieId));
        alert('Movie deleted successfully!');
      } catch (error) {
        console.error('Error deleting movie:', error);
        alert('Failed to delete movie. See the console for more information.');
      }
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      {editingMovie && (
        <EditMovieForm movie={editingMovie} setEditingMovie={setEditingMovie} />
      )}
      <table className="table is-fullwidth">
        <thead>
          <tr>
            <th>Title</th>
            <th>Directors</th>
            <th>Actors</th>
            <th>Genres</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {movies.map((movie) => (
            <tr key={movie.id}>
              <td>{movie.title}</td>
              <td>{movie.directors.join(', ')}</td>
              <td>{movie.actors.join(', ')}</td>
              <td>{movie.genres.join(', ')}</td>
              <td>
                <button
                  type="button"
                  className="button is-info is-small"
                  onClick={() => setEditingMovie(movie)}
                >
                  Edit
                </button>{' '}
                <button
                  type="button"
                  className="button is-danger is-small"
                  onClick={() => handleDelete(movie.id)}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default MovieTable;
