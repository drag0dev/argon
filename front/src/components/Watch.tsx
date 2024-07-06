import React, { useState, useEffect } from 'react';
import ReactPlayer from 'react-player';
import { useParams } from 'react-router-dom';

const Watch = () => {
  const [videoInfo, setVideoInfo] = useState(null);
  const { id } = useParams();

  useEffect(() => {
    // Simulated fetch to get video details based on the ID
    fetchVideoDetails(id);
  }, [id]);

  const fetchVideoDetails = async (id) => {
    const details = {
      id: id,
      title: 'Inception',
      description: 'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.',
      type: 'Movie',
      videoUrl: 'https://www.youtube.com/watch?v=YoHD9XEInc0',
      genres: ['Action', 'Adventure', 'Sci-Fi'],
      actors: ['Leonardo DiCaprio', 'Joseph Gordon-Levitt', 'Ellen Page'],
      directors: ['Christopher Nolan'],
    };
    // Simulating a fetch call
    setVideoInfo(details);
  };

  const handleSubscribe = (itemType, itemName) => {
    console.log(`Subscribed to ${itemType}: ${itemName}`);
    // Implementation for subscription can go here
  };

  if (!videoInfo) {
    return <div className="notification is-info">Loading...</div>;
  }

  return (
    <div className="container mt-5">
      <h1 className="title">{videoInfo.title}</h1>
      <ReactPlayer url={videoInfo.videoUrl} controls={true} className="mb-4" />
      <div className="content">
        <p><strong>Type:</strong> {videoInfo.type}</p>
        <p><strong>Description:</strong> {videoInfo.description}</p>
        <div>
          <strong>Genres:</strong>
          <ul>
            {videoInfo.genres.map(genre => (
              <li key={genre}>
                {genre} <button type="button" className="button is-small is-info" onClick={() => handleSubscribe('Genre', genre)}>Subscribe</button>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <strong>Actors:</strong>
          <ul>
            {videoInfo.actors.map(actor => (
              <li key={actor}>
                {actor} <button type="button" className="button is-small is-info" onClick={() => handleSubscribe('Actor', actor)}>Subscribe</button>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <strong>Directors:</strong>
          <ul>
            {videoInfo.directors.map(director => (
              <li key={director}>
                {director} <button type="button" className="button is-small is-info" onClick={() => handleSubscribe('Director', director)}>Subscribe</button>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
};

export default Watch;
