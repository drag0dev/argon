import React, { useState, useEffect } from 'react';
import ReactPlayer from 'react-player';
import { useParams } from 'react-router-dom';
import ReviewSection from '../components/ReviewSection';
import { fetchAuthSession } from 'aws-amplify/auth';

const API_URL = process.env.API_URL;

const Watch = () => {
  const [videoInfo, setVideoInfo] = useState(null);
  const [videoUrl, setVideoUrl] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const { uuid, seasonId, episodeId } = useParams();

  useEffect(() => {
    fetchVideoDetails(uuid, seasonId, episodeId);
  }, [uuid, seasonId, episodeId]);

  const fetchVideoDetails = async (id, seasonId, episodeId) => {
    setIsLoading(true);
    setError(null);

    try {
      const session = await fetchAuthSession();
      let token = session.tokens?.idToken!.toString();

      let url: string;
      if (seasonId && episodeId) {
        // Fetch TV show episode
        url = `${API_URL}/tvShow?uuid=${uuid}&season=${seasonId}&episode=${episodeId}&resolution=1920:1080`;
      } else {
        // Fetch movie
        url = `${API_URL}/movie?uuid=${uuid}&resolution=1920:1080`;
      }

      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (!response.ok) {
        throw new Error('Failed to fetch video details');
      }

      const data = await response.json();

      if (seasonId && episodeId) {
        const season = data.show.seasons.find(
          (season) => season.seasonNumber === +seasonId,
        );
        const episode = season.episodes.find(
          (episode) => episode.episodeNumber === +episodeId,
        );

        setVideoInfo(episode);
        setVideoUrl(data.url);
      } else {
        setVideoInfo(data.movie);
        setVideoUrl(data.url);
      }
    } catch (error) {
      console.error('Error fetching video details:', error);
      setError('Failed to load video details. Please try again later.');
    } finally {
      setIsLoading(false);
    }
  };

  const SubscriptionType = {
    Actor: 0,
    Director: 1,
    Genre: 2,
  };

  const handleSubscribe = async (type, item) => {
    try {
      const { tokens, identityId } = await fetchAuthSession();
      const userId = tokens.idToken.payload.sub; // 'sub' claim contains the user's UUID
      console.log(`Subscribed to ${type}: ${item} by user ${userId}`);

      const subscriptionData = {
        UserID: userId,
        Type: SubscriptionType[type],
        Target: item,
      };

      const response = await fetch(`${API_URL}/subscription`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokens.idToken.toString()}`,
        },
        body: JSON.stringify(subscriptionData),
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const responseData = await response.json();
      console.log('Subscription successful:', responseData);
      return responseData;
    } catch (error) {
      console.error('Failed to subscribe:', error);
      throw error;
    }
  };

  if (isLoading) {
    return <div className="notification is-info">Loading...</div>;
  }

  if (error) {
    return <div className="notification is-danger">{error}</div>;
  }

  if (!videoInfo) {
    return (
      <div className="notification is-warning">
        No video information available.
      </div>
    );
  }

  return (
    <div className="container mt-5">
      <h1 className="title">{videoInfo.title}</h1>
      {seasonId && episodeId && (
        <h3 className="subtitle">
          Season {seasonId}, Episode {episodeId}
        </h3>
      )}
      <ReactPlayer url={videoUrl} controls={true} className="mb-4" />
      <div className="content">
        <p>
          <strong>Description:</strong> {videoInfo.description}
        </p>
        {!seasonId && !episodeId && (
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
        )}
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
      <ReviewSection targetUUID={uuid} />
    </div>
  );
};

export default Watch;
