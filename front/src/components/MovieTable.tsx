import React, { useState, useEffect } from 'react';
import { API_URL } from '../../config.ts';

const MovieTable = () => {
  const [movies, setMovies] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchMovies = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(`${API_URL}/api/movies`);
        const data = await response.json();
        setMovies(data);
      } catch (error) {
        console.error('Failed to fetch movies:', error);
      }
      setIsLoading(false);
    };

    fetchMovies();
  }, []);

  const handleDelete = async (movieId) => {
    const confirmed = window.confirm('Are you sure you want to delete this movie?');
    if (confirmed) {
      setIsLoading(true);
      try {
        await fetch(`${API_URL}/api/movies/${movieId}`, {
          method: 'DELETE',
        });
        // Remove movie from state after successful deletion
        setMovies(movies.filter(movie => movie.id !== movieId));
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
      <table className="table is-fullwidth">
        <thead>
          <tr>
            <th>Title</th>
            <th>Description</th>
            <th>Genres</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {movies.map((movie) => (
            <tr key={movie.id}>
              <td>{movie.title}</td>
              <td>{movie.description}</td>
              <td>{movie.genres.join(', ')}</td>
              <td>
                <button
                  type="button"
                  className="button is-danger is-small"
                  onClick={() => handleDelete(movie.id)}
                  disabled={isLoading}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {isLoading && <p>Loading...</p>}
    </div>
  );
};

export default MovieTable;
