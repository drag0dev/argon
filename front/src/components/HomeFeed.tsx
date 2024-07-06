import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom'; // Assuming you're using react-router for navigation

const HomeFeed = () => {
  const [recommendations, setRecommendations] = useState([]);
  const [otherRecommendations, setOtherRecommendations] = useState([]); // New state for other shows

  useEffect(() => {
    // Simulated fetch call to get recommendations
    fetchRecommendations();
  }, []);

  useEffect(() => {
    // Simulated fetch call to get other recommendations
    fetchOtherRecommendations();
  }, []);

  const fetchRecommendations = async () => {
    const fetchedData = [
      {
        id: 1,
        title: 'Inception',
        description:
          'A thief who steals corporate secrets through dream-sharing technology...',
        type: 'Movie',
      },
      {
        id: 2,
        title: 'Stranger Things',
        description:
          'When a young boy disappears, his mother, a police chief, and his friends must confront terrifying supernatural forces...',
        type: 'TV Show',
      },
      {
        id: 3,
        title: 'Interstellar',
        description:
          "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival...",
        type: 'Movie',
      },
      {
        id: 4,
        title: 'The Matrix',
        description:
          'A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers...',
        type: 'Movie',
      },
      {
        id: 5,
        title: 'Game of Thrones',
        description:
          'Nine noble families fight for control over the lands of Westeros, while an ancient enemy returns after being dormant for millennia...',
        type: 'TV Show',
      },
      {
        id: 6,
        title: 'The Witcher',
        description:
          'Geralt of Rivia, a solitary monster hunter, struggles to find his place in a world where people often prove more wicked than beasts...',
        type: 'TV Show',
      },
    ];
    setRecommendations(fetchedData);
  };

  const fetchOtherRecommendations = async () => {
    // Simulated fetch for other shows
    const otherData = [
      {
        id: 7,
        title: 'Black Mirror',
        description:
          "An anthology series exploring a twisted, high-tech multiverse where humanity's greatest innovations and darkest instincts collide.",
        type: 'TV Show',
      },
      {
        id: 8,
        title: 'Breaking Bad',
        description:
          'A high school chemistry teacher turned methamphetamine manufacturing drug dealer teams with a former student...',
        type: 'TV Show',
      },
      {
        id: 9,
        title: 'Chernobyl',
        description:
          'A dramatization of the true story of one of the worst man-made catastrophes in history, the catastrophic nuclear accident at Chernobyl.',
        type: 'TV Show',
      },
      {
        id: 10,
        title: 'The Mandalorian',
        description:
          'A lone bounty hunter in the outer reaches of the galaxy, far from the authority of the New Republic...',
        type: 'TV Show',
      },
      {
        id: 11,
        title: 'Arcane',
        description:
          'Set in the utopian region of Piltover and the oppressed underground of Zaun, the story follows the origins of two iconic League champions-and the power that will tear them apart.',
        type: 'TV Show',
      },
      {
        id: 12,
        title: 'Westworld',
        description:
          'Set at the intersection of the near future and the reimagined past, explore a world in which every human appetite can be indulged without consequence.',
        type: 'TV Show',
      },
    ];
    setOtherRecommendations(otherData); // Setting the other recommendations
  };

  return (
    <div className="container">
      <h1 className="title is-3">Recommended for You</h1>
      <div className="grid is-col-min-10">
        {recommendations.map((rec) => (
          <div key={rec.id} className="cell is-flex">
            <div className="card is-flex is-flex-direction-column">
              <header className="card-header">
                <p className="card-header-title">{rec.title}</p>
              </header>
              <div className="card-content is-flex-grow-1">
                {rec.description}
              </div>
              <footer className="card-footer">
                <Link
                  to={`/${rec.type.toLowerCase().replace(/\s/g, '')}/${rec.id}/details`}
                  className="card-footer-item"
                >
                  Watch Now
                </Link>
              </footer>
            </div>
          </div>
        ))}
      </div>
      {/* New section for exploring other shows */}
      <div className="section">
        <h2 className="title is-4">Explore Other Shows</h2>
        <div className="grid is-col-min-10">
          {otherRecommendations.map(
            (
              other, // Added new mapping for other shows
            ) => (
              <div key={other.id} className="cell is-flex">
                <div className="card is-flex is-flex-direction-column">
                  <header className="card-header">
                    <p className="card-header-title">{other.title}</p>
                  </header>
                  <div className="card-content is-flex-grow-1">
                    {other.description}
                  </div>
                  <footer className="card-footer">
                    <Link
                      to={`/${other.type.toLowerCase().replace(/\s/g, '')}/${other.id}/details`}
                      className="card-footer-item"
                    >
                      Watch Now
                    </Link>
                  </footer>
                </div>
              </div>
            ),
          )}
        </div>
      </div>
    </div>
  );
};

export default HomeFeed;
