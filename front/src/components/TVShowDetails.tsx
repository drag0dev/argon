import React from 'react';
import { Link } from 'react-router-dom';
import SeasonDetails from './SeasonDetails';
import { useEffect } from 'react';
import { useState } from 'react';
import { useParams } from 'react-router-dom';

import { fetchAuthSession } from 'aws-amplify/auth';

const API_URL = process.env.API_URL;

const TVShowDetails = () => {
  const [tvShow, setTVShow] = React.useState(null);
  const { uuid } = useParams();

  useEffect(() => {
    fetchTVShowDetails(uuid);
  }, [uuid]);

  const fetchTVShowDetails = async (uuid: string) => {
    try {
      const session = await fetchAuthSession();
      let token = session.tokens?.idToken!.toString();

      const url = `${API_URL}/tvShow?uuid=${uuid}&resolution=1920:1080&season=1&episode=1`;
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch TV show details');
      }

      const data = await response.json();
      setTVShow(data.show);
    } catch (error) {
      console.error('Error fetching TV show details:', error);
    }
  };

  const handleSubscribe = (type, item) => {
    console.log(`Subscribed to ${type}: ${item}`);
    // Subscription logic here
  };

  return (
    <section className="section">
      <div className="container">
        { !tvShow && (
          <div className="notification is-info">
            Loading TV show details...
          </div>
        )}

        { tvShow && (
        <div>
        <h1 className="title">{tvShow.title}</h1>
        <div>
          <strong>Genres:</strong>
          <ul>
            {tvShow.genres.map((genre) => (
              <li key={genre}>
                {genre}
                <button
                  type="button"
                  onClick={() => handleSubscribe('genre', genre)}
                  className="button is-small is-info ml-2"
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
            {tvShow.directors.map((director) => (
              <li key={director}>
                {director}
                <button
                  type="button"
                  onClick={() => handleSubscribe('director', director)}
                  className="button is-small is-info ml-2"
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
            {tvShow.actors.map((actor) => (
              <li key={actor}>
                {actor}
                <button
                  type="button"
                  onClick={() => handleSubscribe('actor', actor)}
                  className="button is-small is-info ml-2"
                >
                  Subscribe
                </button>
              </li>
            ))}
          </ul>
        </div>
        {tvShow.seasons.map((season) => (
          <SeasonDetails
            key={season.seasonNumber}
            season={season}
            showId={tvShow.id}
          />
        ))}
        </div>
        )}
      </div>
    </section>
  );
};

export default TVShowDetails;
