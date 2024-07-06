import React from 'react';
import { Link } from 'react-router-dom';

const EpisodeDetails = ({ title, showId, seasonNum, episodeNum }) => {
  return (
    <li>
      {title} - <Link to={`/tvshow/${showId}/watch/${seasonNum}/${episodeNum}`} className="button is-small is-primary">Watch</Link>
    </li>
  );
};

export default EpisodeDetails;
