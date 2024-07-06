import React from 'react';
import EpisodeDetails from './EpisodeDetails';

const SeasonDetails = ({ season, showId }) => {
  return (
    <div className="container mt-5">
      <div className="box">
        <h2 className="subtitle">Season {season.seasonNumber}</h2>
        <ul>
          {season.episodes.map((episode) => (
            <EpisodeDetails
              key={episode.episodeNumber}
              title={episode.title}
              showId={showId}
              seasonNum={season.seasonNumber}
              episodeNum={episode.episodeNumber}
            />
          ))}
        </ul>
      </div>
    </div>
  );
};

export default SeasonDetails;
