import MovieTable from '../components/MovieTable';
import React, { useState } from 'react';
import { fetchAuthSession } from 'aws-amplify/auth';

const API_URL = process.env.API_URL;
const EditOrDeleteMoviePage = () => {
  const [uuid, setUuid] = useState('');
  const [seasonNum, setSeasonNum] = useState('');
  const [episodeNum, setEpisodeNum] = useState('');

  const handleDelete = async () => {
    try {

      const session = await fetchAuthSession();
      let token  = session.tokens?.idToken!.toString()

      let url = `${API_URL}/tvShow?uuid=${uuid}`;

      if (seasonNum != '') {
        url += `&season=${seasonNum}`;
      }

      if (episodeNum != '') {
        url += `&episode=${episodeNum}`;
      }

      await fetch(url, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
      });
      alert('Movie deleted successfully');
    } catch (error) {
      console.error('Error deleting movie:', error);
      alert('Failed to delete movie');
    }
  };

  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Edit/Delete Movie</h1>
        <div className="field">
          <label className="label">Show UUID</label>
          <div className="control">
            <input 
              className="input" 
              type="text" 
              value={uuid} 
              onChange={(e) => setUuid(e.target.value)} 
              placeholder="Enter show UUID" 
            />
          </div>
        </div>
        <div className="field">
          <label className="label">Season</label>
          <div className="control">
            <input 
              className="input" 
              type="text" 
              value={seasonNum} 
              onChange={(e) => setSeasonNum(e.target.value)} 
              placeholder="Enter season num" 
            />
          </div>
        </div>
        <div className="field">
          <label className="label">episdoe</label>
          <div className="control">
            <input 
              className="input" 
              type="text" 
              value={episodeNum} 
              onChange={(e) => setEpisodeNum(e.target.value)} 
              placeholder="Enter episode num" 
            />
          </div>
        </div>
        <div className="buttons">
          <button className="button is-danger" onClick={handleDelete}>Delete Movie</button>
        </div>
        <MovieTable />
      </div>
    </section>
  );
}

export default EditOrDeleteMoviePage;

