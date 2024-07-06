import React, { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSave } from '@fortawesome/free-solid-svg-icons';

const API_URL = process.env.API_URL;

interface TVShow {
  id: string;
  title: string;
  genres: string[];
  actors: string[];
  directors: string[];
}

const EditTVShowForm = ({ tvShow: initialTVShow, setEditingTVShow }) => {
  const [tvShow, setTVShow] = useState<TVShow>(initialTVShow);
  const [isLoading, setIsLoading] = useState(false);

  // Update state when initialTVShow changes
  useEffect(() => {
    setTVShow(initialTVShow);
  }, [initialTVShow]);  // Dependency on initialTVShow

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setTVShow({ ...tvShow, [name]: value });
  };

  const handleArrayInputChange = (e, field) => {
    const { value } = e.target;
    setTVShow({
      ...tvShow,
      [field]: value.split(',').map(item => item.trim()),
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const tvShowData = { ...tvShow }; // Prepare tv show data
      const url = `${API_URL}/api/tvshows/${tvShow.id}`; // Adjust URL for update

      const response = await fetch(url, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(tvShowData),
      });

      if (!response.ok) {
        throw new Error('Failed to update tv show');
      }

      alert('TV show updated successfully!');
      setEditingTVShow(null); // Close the form upon success
    } catch (error) {
      console.error('Error updating tv show:', error);
      alert('Failed to update tv show. See the console for more information.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Edit TV Show Information</h2>
        
        <div className="field">
          <label className="label" htmlFor="title">Title</label>
          <div className="control">
            <input className="input" type="text" id="title" name="title" value={tvShow.title} onChange={handleInputChange} required />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="genres">Genres (comma-separated)</label>
          <div className="control">
            <input className="input" type="text" id="genres" name="genres" value={tvShow.genres.join(', ')} onChange={(e) => handleArrayInputChange(e, 'genres')} />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="actors">Actors (comma-separated)</label>
          <div className="control">
            <input className="input" type="text" id="actors" name="actors" value={tvShow.actors.join(', ')} onChange={(e) => handleArrayInputChange(e, 'actors')} />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="directors">Directors (comma-separated)</label>
          <div className="control">
            <input className="input" type="text" id="directors" name="directors" value={tvShow.directors.join(', ')} onChange={(e) => handleArrayInputChange(e, 'directors')} />
          </div>
        </div>

        <div className="field">
          <div className="control">
            <button type="submit" className="button is-primary" disabled={isLoading}>
              <span className="icon"><FontAwesomeIcon icon={faSave} /></span>
              <span>{isLoading ? 'Saving...' : 'Save Changes'}</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default EditTVShowForm;
