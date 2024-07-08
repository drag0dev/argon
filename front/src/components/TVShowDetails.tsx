import React from 'react';
import { Link } from 'react-router-dom';
import SeasonDetails from './SeasonDetails';
import { useEffect } from 'react';
import { useState } from 'react';
import { useParams } from 'react-router-dom';

import { fetchAuthSession } from 'aws-amplify/auth';

const dummyDetails = {
  id: 1,
  title: 'placeholder data ? (stuff didnt load)',
  genres: ['Drama', 'Fantasy', 'Horror'],
  actors: ['Winona Ryder', 'David Harbour', 'Finn Wolfhard'],
  directors: ['The Duffer Brothers'],
  seasons: [
    {
      seasonNumber: 1,
      episodes: [
        {
          episodeNumber: 1,
          title: 'The Vanishing of Will Byers',
          description:
            'A young boy disappears, leading to an investigation involving supernatural forces.',
          actors: ['Winona Ryder', 'David Harbour', 'Finn Wolfhard'],
          directors: ['The Duffer Brothers'],
        },
        {
          episodeNumber: 2,
          title: 'The Weirdo on Maple Street',
          description:
            "A girl with a shaved head and strange powers appears, providing a clue to Will's disappearance.",
          actors: ['Winona Ryder', 'David Harbour', 'Finn Wolfhard'],
          directors: ['The Duffer Brothers'],
        },
      ],
    },
    {
      seasonNumber: 2,
      episodes: [
        {
          episodeNumber: 1,
          title: 'MADMAX',
          description:
            'The boys encounter a new girl at school while supernatural events continue to plague the town.',
          actors: ['Winona Ryder', 'David Harbour', 'Finn Wolfhard'],
          directors: ['The Duffer Brothers'],
        },
        {
          episodeNumber: 2,
          title: 'Trick or Treat, Freak',
          description:
            'Will struggles to adjust to life after the Upside Down as Halloween approaches.',
          actors: ['Winona Ryder', 'David Harbour', 'Finn Wolfhard'],
          directors: ['The Duffer Brothers'],
        },
      ],
    },
  ],
};

const TVShowDetails = () => {
  const [tvShow, setTVShow] = React.useState(dummyDetails);
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
    </section>
  );
};

export default TVShowDetails;
