import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faMinus } from '@fortawesome/free-solid-svg-icons';
import VideoUpload from './VideoUpload';

const API_URL = process.env.API_URL;

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
    updatedEpisodes[episodeIndex] = { ...updatedEpisodes[episodeIndex], [field]: value };
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
      seasons: [...show.seasons, { seasonNumber: show.seasons.length + 1, episodes: [] }],
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
    updatedSeasons[seasonIndex].episodes = updatedSeasons[seasonIndex].episodes.filter(
      (_, i) => i !== episodeIndex
    );
    setShow({ ...show, seasons: updatedSeasons });
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
                        handleEpisodeChange(seasonIndex, episodeIndex, 'title', e.target.value)
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
                        handleEpisodeChange(seasonIndex, episodeIndex, 'description', e.target.value)
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
                          e.target.value.split(',').map((item) => item.trim())
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
                          e.target.value.split(',').map((item) => item.trim())
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

        <button type="button" className="button is-primary mt-4" onClick={addSeason}>
          <FontAwesomeIcon icon={faPlus} /> Add Season
        </button>

        <div className="field mt-6">
          <div className="control">
            <button
              type="submit"
              className="button is-primary"
              disabled={isLoading || show.seasons.some((s) => s.episodes.length === 0)}
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
