import React, { useState, useEffect } from 'react';
import EditSeasonForm from './EditSeasonForm'; // Assuming this is the correct path
import EditEpisodeForm from './EditEpisodeForm'; // Assuming this is the correct path
import EditTVShowForm from './EditTVShowForm'; // Assuming this is the correct path

const dummyTVShows = [
  {
    id: 1,
    title: 'Stranger Things',
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
  },
  {
    id: 2,
    title: 'Breaking Bad',
    genres: ['Crime', 'Drama', 'Thriller'],
    actors: ['Bryan Cranston', 'Aaron Paul', 'Anna Gunn'],
    directors: ['Vince Gilligan'],
    seasons: [
      {
        seasonNumber: 1,
        episodes: [
          {
            episodeNumber: 1,
            title: 'Pilot',
            description:
              'A high school chemistry teacher turns to making and selling methamphetamine.',
            actors: ['Bryan Cranston', 'Aaron Paul', 'Anna Gunn'],
            directors: ['Vince Gilligan'],
          },
          {
            episodeNumber: 2,
            title: "Cat's in the Bag...",
            description:
              'Walter and Jesse attempt to dispose of the bodies from their first cook.',
            actors: ['Bryan Cranston', 'Aaron Paul', 'Anna Gunn'],
            directors: ['Vince Gilligan'],
          },
        ],
      },
      {
        seasonNumber: 2,
        episodes: [
          {
            episodeNumber: 1,
            title: 'Seven Thirty-Seven',
            description:
              "Walter and Jesse's operation faces new threats and challenges.",
            actors: ['Bryan Cranston', 'Aaron Paul', 'Anna Gunn'],
            directors: ['Vince Gilligan'],
          },
          {
            episodeNumber: 2,
            title: 'Grilled',
            description:
              'Tuco takes Walter and Jesse hostage as the DEA closes in.',
            actors: ['Bryan Cranston', 'Aaron Paul', 'Anna Gunn'],
            directors: ['Vince Gilligan'],
          },
        ],
      },
    ],
  },
  {
    id: 3,
    title: 'The Witcher',
    genres: ['Action', 'Adventure', 'Fantasy'],
    actors: ['Henry Cavill', 'Anya Chalotra', 'Freya Allan'],
    directors: ['Lauren Schmidt Hissrich'],
    seasons: [
      {
        seasonNumber: 1,
        episodes: [
          {
            episodeNumber: 1,
            title: "The End's Beginning",
            description:
              'Geralt of Rivia, a mutated monster hunter, struggles to find his place in a world where people often prove more wicked than beasts.',
            actors: ['Henry Cavill', 'Anya Chalotra', 'Freya Allan'],
            directors: ['Lauren Schmidt Hissrich'],
          },
          {
            episodeNumber: 2,
            title: 'Four Marks',
            description:
              "Yennefer's early days as a sorceress and her path to power are revealed.",
            actors: ['Henry Cavill', 'Anya Chalotra', 'Freya Allan'],
            directors: ['Lauren Schmidt Hissrich'],
          },
        ],
      },
      {
        seasonNumber: 2,
        episodes: [
          {
            episodeNumber: 1,
            title: 'A Grain of Truth',
            description:
              'Geralt reunites with an old friend as he seeks safety for Ciri.',
            actors: ['Henry Cavill', 'Anya Chalotra', 'Freya Allan'],
            directors: ['Lauren Schmidt Hissrich'],
          },
          {
            episodeNumber: 2,
            title: 'Kaer Morhen',
            description: 'Ciri trains with the witchers at their fortress.',
            actors: ['Henry Cavill', 'Anya Chalotra', 'Freya Allan'],
            directors: ['Lauren Schmidt Hissrich'],
          },
        ],
      },
    ],
  },
  {
    id: 4,
    title: 'Black Mirror',
    genres: ['Drama', 'Sci-Fi', 'Thriller'],
    actors: ['Bryce Dallas Howard', 'Daniel Kaluuya', 'Jon Hamm'],
    directors: ['Charlie Brooker'],
    seasons: [
      {
        seasonNumber: 1,
        episodes: [
          {
            episodeNumber: 1,
            title: 'The National Anthem',
            description:
              'A twisted tale of a prime minister faced with a horrifying choice.',
            actors: ['Bryce Dallas Howard', 'Daniel Kaluuya', 'Jon Hamm'],
            directors: ['Charlie Brooker'],
          },
          {
            episodeNumber: 2,
            title: 'Fifteen Million Merits',
            description:
              'In a dystopian future, society is controlled by technology and the meritocracy.',
            actors: ['Bryce Dallas Howard', 'Daniel Kaluuya', 'Jon Hamm'],
            directors: ['Charlie Brooker'],
          },
        ],
      },
      {
        seasonNumber: 2,
        episodes: [
          {
            episodeNumber: 1,
            title: 'Be Right Back',
            description:
              'A grieving woman uses technology to reconnect with her deceased partner.',
            actors: ['Bryce Dallas Howard', 'Daniel Kaluuya', 'Jon Hamm'],
            directors: ['Charlie Brooker'],
          },
          {
            episodeNumber: 2,
            title: 'White Bear',
            description:
              'A woman awakes in a strange dystopian world where she is relentlessly pursued.',
            actors: ['Bryce Dallas Howard', 'Daniel Kaluuya', 'Jon Hamm'],
            directors: ['Charlie Brooker'],
          },
        ],
      },
    ],
  },
];

const TVShowTable = () => {
  const [tvShows, setTVShows] = useState(dummyTVShows);
  const [openShowId, setOpenShowId] = useState(null);
  const [openSeasonNumber, setOpenSeasonNumber] = useState(null);
  const [editingSeason, setEditingSeason] = useState(null);
  const [editingEpisode, setEditingEpisode] = useState(null);
  const [editingTVShow, setEditingTVShow] = useState(null);

  const handleEditTVShow = (showId) => {
    const tvShow = tvShows.find((show) => show.id === showId);
    if (tvShow) {
      setEditingTVShow(tvShow);
    }
  };

  // Toggle the visibility of show details
  const toggleShowDetails = (showId) => {
    setOpenShowId(openShowId === showId ? null : showId);
    setOpenSeasonNumber(null); // Reset the open season when toggling different shows
    setEditingSeason(null);
    setEditingEpisode(null);
  };

  // Toggle the visibility of season details
  const toggleSeasonDetails = (seasonNumber) => {
    setOpenSeasonNumber(
      openSeasonNumber === seasonNumber ? null : seasonNumber,
    );
    setEditingEpisode(null);
  };

  // Function to start editing a season
  const handleEditSeason = (showId, seasonNumber) => {
    const season = tvShows
      .find((show) => show.id === showId)
      ?.seasons.find((season) => season.seasonNumber === seasonNumber);
    if (season) {
      setEditingSeason({ showId, ...season });
    }
  };

  // Function to start editing an episode
  const handleEditEpisode = (showId, seasonNumber, episodeNumber) => {
    const episode = tvShows
      .find((show) => show.id === showId)
      ?.seasons.find((season) => season.seasonNumber === seasonNumber)
      ?.episodes.find((episode) => episode.episodeNumber === episodeNumber);
    if (episode) {
      setEditingEpisode({ showId, seasonNumber, ...episode });
    }
  };

  // Placeholder for deleting an episode
  const handleDeleteEpisode = (showId, seasonNumber, episodeNumber) => {
    console.log(
      `Delete Episode ${episodeNumber} of Season ${seasonNumber} of Show ${showId}`,
    );
    // Actual delete logic here
  };

  return (
    <div className="container">
      {editingSeason && (
        <EditSeasonForm
          season={editingSeason}
          setEditingSeason={setEditingSeason}
        />
      )}
      {editingEpisode && (
        <EditEpisodeForm
          episode={editingEpisode}
          setEditingEpisode={setEditingEpisode}
        />
      )}
      {editingTVShow && (
        <EditTVShowForm
          tvShow={editingTVShow}
          setEditingTVShow={setEditingTVShow}
        />
      )}
      <table className="table is-fullwidth">
        <thead>
          <tr>
            <th>Title</th>
            <th>Seasons</th>
            <th>Episodes</th>
            <th>Genres</th>
            <th>Actors</th>
            <th>Directors</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {tvShows.map((show) => (
            <tr key={show.id}>
              <td>{show.title}</td>
              <td>{show.seasons.length}</td>
              <td>
                {show.seasons.reduce(
                  (acc, season) => acc + season.episodes.length,
                  0,
                )}
              </td>
              <td>{show.genres.join(', ')}</td>
              <td>{show.actors.join(', ')}</td>
              <td>{show.directors.join(', ')}</td>
              <td>
                <button
                  type="button"
                  className="button is-small"
                  onClick={() => toggleShowDetails(show.id)}
                >
                  {openShowId === show.id ? 'Close' : 'More'}
                </button>
                <button
                  type="button"
                  className="button is-small"
                  onClick={() => handleEditTVShow(show.id)}
                >
                  Edit
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {openShowId && (
        <div className="box">
          {tvShows
            .find((show) => show.id === openShowId)
            .seasons.map((season) => (
              <div key={season.seasonNumber}>
                <h5>Season {season.seasonNumber}</h5>
                <button
                  type="button"
                  className="button is-small"
                  onClick={() =>
                    handleEditSeason(openShowId, season.seasonNumber)
                  }
                >
                  Edit
                </button>
                <button
                  type="button"
                  className="button is-small"
                  onClick={() => toggleSeasonDetails(season.seasonNumber)}
                >
                  {openSeasonNumber === season.seasonNumber ? 'Close' : 'More'}
                </button>

                {openSeasonNumber === season.seasonNumber && (
                  <div className="box">
                    {season.episodes.map((episode) => (
                      <div key={episode.episodeNumber}>
                        <p>
                          Episode {episode.episodeNumber}: {episode.title}
                        </p>
                        <button
                          type="button"
                          className="button is-small"
                          onClick={() =>
                            handleEditEpisode(
                              openShowId,
                              season.seasonNumber,
                              episode.episodeNumber,
                            )
                          }
                        >
                          Edit
                        </button>
                        <button
                          type="button"
                          className="button is-small"
                          onClick={() =>
                            handleDeleteEpisode(
                              openShowId,
                              season.seasonNumber,
                              episode.episodeNumber,
                            )
                          }
                        >
                          Delete
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
        </div>
      )}
    </div>
  );
};

export default TVShowTable;
