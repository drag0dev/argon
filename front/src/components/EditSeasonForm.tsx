import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEdit } from '@fortawesome/free-solid-svg-icons';
const API_URL = process.env.API_URL;

interface Season {
  id: string;
  seasonNumber: number;
  description: string;
}

const EditSeasonForm = ({ season: initialSeason, setEditingSeason }) => {
  const [season, setSeason] = useState<Season>(initialSeason);
  const [isLoading, setIsLoading] = useState(false);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setSeason({ ...season, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const seasonData = { ...season };
      const url = `${API_URL}/api/seasons/${season.id}`;

      const response = await fetch(url, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(seasonData),
      });

      if (!response.ok) {
        throw new Error('Failed to update season');
      }

      alert('Season updated successfully!');
      setEditingSeason(null);
    } catch (error) {
      console.error('Error updating season:', error);
      alert('Failed to update season. See the console for more information.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Edit Season Information</h2>
        <div className="field">
          <label className="label" htmlFor="seasonNumber">Season Number</label>
          <div className="control">
            <input className="input" type="number" id="seasonNumber" name="seasonNumber" value={season.seasonNumber} onChange={handleInputChange} required />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="description">Description (optional)</label>
          <div className="control">
            <textarea className="textarea" id="description" name="description" value={season.description} onChange={handleInputChange} />
          </div>
        </div>

        <div className="field">
          <div className="control">
            <button type="submit" className="button is-primary" disabled={isLoading}>
              <span className="icon"><FontAwesomeIcon icon={faEdit} /></span>
              <span>{isLoading ? 'Updating...' : 'Update Season'}</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default EditSeasonForm;
