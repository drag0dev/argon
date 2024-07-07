import React, { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faMinus } from '@fortawesome/free-solid-svg-icons';
import VideoUpload from './VideoUpload';

const API_URL = process.env.API_URL;

const dummyTVShows = [
  {
    id: 1,
    title: 'Stranger Things',
    genres: ['Drama', 'Fantasy', 'Horror'],
    actors: ['Winona Ryder', 'David Harbour', 'Finn Wolfhard'],
    directors: ['The Duffer Brothers'],
    description:
      "In a small town where supernatural forces loom, a young boy's mysterious disappearance sets off a chain of events uncovering government experiments and alternate dimensions.",
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
    description:
      "A high school chemistry teacher, diagnosed with terminal cancer, turns to manufacturing and selling methamphetamine to secure his family's financial future, leading to a dangerous descent into the criminal underworld.",
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
    description:
      'Geralt of Rivia, a solitary monster hunter, navigates a world where powerful sorceresses, cunning kings, and dangerous creatures vie for dominance, while destiny binds him to a young princess with a mysterious past.',
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
    description:
      'Each episode of this anthology series explores a standalone story, often dystopian and thought-provoking, highlighting the dark side of technology and its impact on modern society through twisted, provocative narratives.',
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

const AddTVShowForm = () => {
  const [show, setShow] = useState({
    title: '',
    description: '',
    genres: [],
    actors: [],
    directors: [],
    seasons: [{ seasonNumber: 1, episodes: [] }],
  });

  const [isLoading, setIsLoading] = useState(false);

  // Function to insert dummy data
  const insertDummyData = (n) => {
    const selectedShow = dummyTVShows[n];
    setShow({
      title: selectedShow.title,
      description: selectedShow.description || '',
      genres: selectedShow.genres,
      actors: selectedShow.actors,
      directors: selectedShow.directors,
      seasons: selectedShow.seasons.map((season) => ({
        seasonNumber: season.seasonNumber,
        episodes: season.episodes.map((episode) => ({
          ...episode,
          video: null, // Set video to null as we don't have actual video files for dummy data
        })),
      })),
    });
  };

  useEffect(() => {
    window.insertDummyData = insertDummyData;
    return () => {
      window.insertDummyData = undefined;
    };
  }, []);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setShow({ ...show, [name]: value });
  };

  const handleArrayInputChange = (e, field) => {
    const { value } = e.target;
    setShow({
      ...show,
      [field]: value.split(',').map((item) => item.trim()),
    });
  };

  const handleSeasonChange = (index, field, value) => {
    const updatedSeasons = [...show.seasons];
    updatedSeasons[index] = { ...updatedSeasons[index], [field]: value };
    setShow({ ...show, seasons: updatedSeasons });
  };

  const handleEpisodeChange = (seasonIndex, episodeIndex, field, value) => {
    const updatedSeasons = [...show.seasons];
    const updatedEpisodes = [...updatedSeasons[seasonIndex].episodes];
    updatedEpisodes[episodeIndex] = {
      ...updatedEpisodes[episodeIndex],
      [field]: value,
    };
    updatedSeasons[seasonIndex].episodes = updatedEpisodes;
    setShow({ ...show, seasons: updatedSeasons });
  };

  const handleVideoUpload = (seasonIndex, episodeIndex, file, metadata) => {
    const updatedSeasons = [...show.seasons];
    const updatedEpisodes = [...updatedSeasons[seasonIndex].episodes];
    updatedEpisodes[episodeIndex].video = {
      file,
      fileType: metadata.type,
      fileSize: metadata.size,
      creationTimestamp: metadata.lastModified.getTime(),
      lastChangeTimestamp: Date.now(),
    };
    updatedSeasons[seasonIndex].episodes = updatedEpisodes;
    setShow({ ...show, seasons: updatedSeasons });
  };

  const addSeason = () => {
    setShow({
      ...show,
      seasons: [
        ...show.seasons,
        { seasonNumber: show.seasons.length + 1, episodes: [] },
      ],
    });
  };

  const removeSeason = (index) => {
    const updatedSeasons = show.seasons.filter((_, i) => i !== index);
    setShow({ ...show, seasons: updatedSeasons });
  };

  const addEpisode = (seasonIndex) => {
    const updatedSeasons = [...show.seasons];
    updatedSeasons[seasonIndex].episodes.push({
      episodeNumber: updatedSeasons[seasonIndex].episodes.length + 1,
      title: '',
      description: '',
      actors: [],
      directors: [],
      video: null,
    });
    setShow({ ...show, seasons: updatedSeasons });
  };

  const removeEpisode = (seasonIndex, episodeIndex) => {
    const updatedSeasons = [...show.seasons];
    updatedSeasons[seasonIndex].episodes = updatedSeasons[
      seasonIndex
    ].episodes.filter((_, i) => i !== episodeIndex);
    setShow({ ...show, seasons: updatedSeasons });
  };

  const allEpisodesHaveVideos = () => {
    return show.seasons.every((season) =>
      season.episodes.every((episode) => episode.video?.file),
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      // Prepare the show data
      const showData = {
        title: show.title,
        description: show.description,
        genres: show.genres,
        actors: show.actors,
        directors: show.directors,
        seasons: show.seasons.map((season) => ({
          seasonNumber: season.seasonNumber,
          episodes: season.episodes.map((episode) => ({
            episodeNumber: episode.episodeNumber,
            title: episode.title,
            description: episode.description,
            actors: episode.actors,
            directors: episode.directors,
            video: {
              fileType: episode.video.fileType,
              fileSize: episode.video.fileSize,
              creationTimestamp: episode.video.creationTimestamp,
              lastChangeTimestamp: episode.video.lastChangeTimestamp,
            },
          })),
        })),
      };

      const url = `${API_URL}/api/shows`;

      // Step 1: Send show metadata and get upload URLs
      const metadataResponse = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(showData),
      });

      if (!metadataResponse.ok) {
        throw new Error('Failed to submit show metadata');
      }

      const { uploadUrls, showId } = await metadataResponse.json();

      // Step 2: Upload video files
      for (let i = 0; i < show.seasons.length; i++) {
        for (let j = 0; j < show.seasons[i].episodes.length; j++) {
          const episode = show.seasons[i].episodes[j];
          const uploadUrl = uploadUrls[i][j];

          const uploadResponse = await fetch(uploadUrl, {
            method: 'PUT',
            body: episode.video.file,
            headers: {
              'Content-Type': episode.video.fileType,
            },
          });

          if (!uploadResponse.ok) {
            throw new Error(`Failed to upload video file for S${i + 1}E${j + 1}`);
          }
        }
      }

      console.log('TV Show and videos uploaded successfully. Show ID:', showId);

      // Clear the form
      setShow({
        title: '',
        description: '',
        genres: [],
        actors: [],
        directors: [],
        seasons: [{ seasonNumber: 1, episodes: [] }],
      });

      alert('TV Show added successfully!');
    } catch (error) {
      console.error('Error adding TV Show:', error);
      alert('Failed to add TV Show. See the console for more information.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">TV Show Information</h2>

        <div className="field">
          <label className="label" htmlFor="title">
            Title
          </label>
          <div className="control">
            <input
              className="input"
              type="text"
              id="title"
              name="title"
              value={show.title}
              onChange={handleInputChange}
              required
            />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="description">
            Description
          </label>
          <div className="control">
            <textarea
              className="textarea"
              id="description"
              name="description"
              value={show.description}
              onChange={handleInputChange}
              required
            />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="genres">
            Genres (comma-separated)
          </label>
          <div className="control">
            <input
              className="input"
              type="text"
              id="genres"
              name="genres"
              value={show.genres.join(', ')}
              onChange={(e) => handleArrayInputChange(e, 'genres')}
            />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="actors">
            Actors (comma-separated)
          </label>
          <div className="control">
            <input
              className="input"
              type="text"
              id="actors"
              name="actors"
              value={show.actors.join(', ')}
              onChange={(e) => handleArrayInputChange(e, 'actors')}
            />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="directors">
            Directors (comma-separated)
          </label>
          <div className="control">
            <input
              className="input"
              type="text"
              id="directors"
              name="directors"
              value={show.directors.join(', ')}
              onChange={(e) => handleArrayInputChange(e, 'directors')}
            />
          </div>
        </div>

        <h3 className="title is-5">Seasons</h3>
        {show.seasons.map((season, seasonIndex) => (
          <div key={season.seasonNumber} className="box">
            <h4 className="title is-6">Season {season.seasonNumber}</h4>
            <button
              type="button"
              className="button is-small is-danger"
              onClick={() => removeSeason(seasonIndex)}
            >
              <FontAwesomeIcon icon={faMinus} /> Remove Season
            </button>

            <h5 className="title is-6 mt-4">Episodes</h5>
            {season.episodes.map((episode, episodeIndex) => (
              <div key={`${seasonIndex}-${episodeIndex}`} className="box">
                <h6 className="title is-6">Episode {episode.episodeNumber}</h6>
                <button
                  type="button"
                  className="button is-small is-danger"
                  onClick={() => removeEpisode(seasonIndex, episodeIndex)}
                >
                  <FontAwesomeIcon icon={faMinus} /> Remove Episode
                </button>

                <div className="field">
                  <label className="label">Title</label>
                  <div className="control">
                    <input
                      className="input"
                      type="text"
                      value={episode.title}
                      onChange={(e) =>
                        handleEpisodeChange(
                          seasonIndex,
                          episodeIndex,
                          'title',
                          e.target.value,
                        )
                      }
                      required
                    />
                  </div>
                </div>

                <div className="field">
                  <label className="label">Description</label>
                  <div className="control">
                    <textarea
                      className="textarea"
                      value={episode.description}
                      onChange={(e) =>
                        handleEpisodeChange(
                          seasonIndex,
                          episodeIndex,
                          'description',
                          e.target.value,
                        )
                      }
                      required
                    />
                  </div>
                </div>

                <div className="field">
                  <label className="label">Actors (comma-separated)</label>
                  <div className="control">
                    <input
                      className="input"
                      type="text"
                      value={episode.actors.join(', ')}
                      onChange={(e) =>
                        handleEpisodeChange(
                          seasonIndex,
                          episodeIndex,
                          'actors',
                          e.target.value.split(',').map((item) => item.trim()),
                        )
                      }
                    />
                  </div>
                </div>

                <div className="field">
                  <label className="label">Directors (comma-separated)</label>
                  <div className="control">
                    <input
                      className="input"
                      type="text"
                      value={episode.directors.join(', ')}
                      onChange={(e) =>
                        handleEpisodeChange(
                          seasonIndex,
                          episodeIndex,
                          'directors',
                          e.target.value.split(',').map((item) => item.trim()),
                        )
                      }
                    />
                  </div>
                </div>

                <VideoUpload
                  onFileUpload={(file, metadata) =>
                    handleVideoUpload(seasonIndex, episodeIndex, file, metadata)
                  }
                />
              </div>
            ))}

            <button
              type="button"
              className="button is-small is-primary mt-4"
              onClick={() => addEpisode(seasonIndex)}
            >
              <FontAwesomeIcon icon={faPlus} /> Add Episode
            </button>
          </div>
        ))}

        <button
          type="button"
          className="button is-primary mt-4"
          onClick={addSeason}
        >
          <FontAwesomeIcon icon={faPlus} /> Add Season
        </button>

        <div className="field mt-6">
          <div className="control">
            <button
              type="submit"
              className="button is-primary"
              disabled={
                isLoading ||
                show.seasons.some((s) => s.episodes.length === 0) ||
                !allEpisodesHaveVideos()
              }
            >
              <span className="icon">
                <FontAwesomeIcon icon={faPlus} />
              </span>
              <span>{isLoading ? 'Uploading...' : 'Add TV Show'}</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default AddTVShowForm;
