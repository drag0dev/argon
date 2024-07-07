import React, { useState, useEffect } from 'react';
import ReactPlayer from 'react-player';
import { useParams } from 'react-router-dom';
import ReviewSection from '../components/ReviewSection';

const Watch = () => {
  const [videoInfo, setVideoInfo] = useState(null);
  const { id, seasonId, episodeId } = useParams();

  useEffect(() => {
    fetchVideoDetails(id, seasonId, episodeId);
  }, [id, seasonId, episodeId]);

  const fetchVideoDetails = async (id, seasonNum, episodeNum) => {
    let details = {};
    if (seasonId && episodeId) {
      // Dummy data for a TV show episode
      details = {
        id: id,
        title: 'Example Show',
        description: 'Description of a specific episode of a TV show.',
        type: 'Episode',
        videoUrl: 'https://www.example.com/episode.mp4',
        genres: ['Drama'],
        actors: ['Actor A', 'Actor B'],
        directors: ['Director X'],
      };
    } else {
      // Data for a movie
      details = {
        id: id,
        title: 'Inception',
        description:
          'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.',
        type: 'Movie',
        videoUrl: 'https://www.youtube.com/watch?v=YoHD9XEInc0',
        genres: ['Action', 'Adventure', 'Sci-Fi'],
        actors: ['Leonardo DiCaprio', 'Joseph Gordon-Levitt', 'Ellen Page'],
        directors: ['Christopher Nolan'],
      };
    }
    setVideoInfo(details);
  };

  const handleSubscribe = (itemType, itemName) => {
    console.log(`Subscribed to ${itemType}: ${itemName}`);
  };

  if (!videoInfo) {
    return <div className="notification is-info">Loading...</div>;
  }

  return (
    <div className="container mt-5">
      <h1 className="title">{videoInfo.title}</h1>
      {seasonId && episodeId && (
        <h3 className="subtitle">
          Season {seasonId}, Episode {episodeId}
        </h3>
      )}
      <ReactPlayer url={videoInfo.videoUrl} controls={true} className="mb-4" />
      <div className="content">
        <p>
          <strong>Type:</strong> {videoInfo.type}
        </p>
        <p>
          <strong>Description:</strong> {videoInfo.description}
        </p>
        <div>
          <strong>Genres:</strong>
          <ul>
            {videoInfo.genres.map((genre) => (
              <li key={genre}>
                {genre}{' '}
                <button
                  type="button"
                  className="button is-small is-info"
                  onClick={() => handleSubscribe('Genre', genre)}
                >
                  Subscribe
                </button>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <strong>Actors:</strong>
          <ul>
            {videoInfo.actors.map((actor) => (
              <li key={actor}>
                {actor}{' '}
                <button
                  type="button"
                  className="button is-small is-info"
                  onClick={() => handleSubscribe('Actor', actor)}
                >
                  Subscribe
                </button>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <strong>Directors:</strong>
          <ul>
            {videoInfo.directors.map((director) => (
              <li key={director}>
                {director}{' '}
                <button
                  type="button"
                  className="button is-small is-info"
                  onClick={() => handleSubscribe('Director', director)}
                >
                  Subscribe
                </button>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <ReviewSection targetUUID={id} />
    </div>
  );
};

export default Watch;
