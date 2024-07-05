import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEdit } from '@fortawesome/free-solid-svg-icons';
import VideoUpload from './VideoUpload';
import { API_URL } from '../../config.ts';

interface VideoMetadata {
  type: string;
  size: number;
  lastModified: Date;
}

interface Episode {
  id: string;
  episodeNumber: number;
  title: string;
  description: string;
  actors: string[];
  directors: string[];
  video?: {
    file: File;
    fileType: string;
    fileSize: number;
    creationTimestamp: number;
    lastChangeTimestamp: number;
  };
}

const EditEpisodeForm = ({ episode: initialEpisode, setEditingEpisode }) => {
  const [episode, setEpisode] = useState<Episode>({
    ...initialEpisode,
    actors: initialEpisode.actors || [], // Ensure actors is always an array
    directors: initialEpisode.directors || [], // Ensure directors is always an array
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEpisode({ ...episode, [name]: value });
  };

  const handleArrayInputChange = (e, field) => {
    const { value } = e.target;
    setEpisode({
      ...episode,
      [field]: value.split(',').map((item) => item.trim()),
    });
  };

  const handleVideoUpload = (file, metadata: VideoMetadata) => {
    setEpisode({
      ...episode,
      video: {
        file,
        fileType: metadata.type,
        fileSize: metadata.size,
        creationTimestamp: metadata.lastModified.getTime(),
        lastChangeTimestamp: Date.now(),
      },
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const episodeData = { ...episode }; // Prepare episode data
      const url = `${API_URL}/api/episodes/${episode.id}`; // Adjust URL for update

      const response = await fetch(url, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(episodeData),
      });

      if (!response.ok) {
        throw new Error('Failed to update episode');
      }

      alert('Episode updated successfully!');
      setEditingEpisode(null); // Close the form upon success
    } catch (error) {
      console.error('Error updating episode:', error);
      alert('Failed to update episode. See the console for more information.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Edit Episode Information</h2>

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
              value={episode.title}
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
              value={episode.description}
              onChange={handleInputChange}
              required
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
              value={episode.actors.join(', ')}
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
              value={episode.directors.join(', ')}
              onChange={(e) => handleArrayInputChange(e, 'directors')}
            />
          </div>
        </div>

        <VideoUpload onFileUpload={handleVideoUpload} editing={true} />

        <div className="field">
          <div className="control">
            <button
              type="submit"
              className="button is-primary"
              disabled={isLoading}
            >
              <span className="icon">
                <FontAwesomeIcon icon={faEdit} />
              </span>
              <span>{isLoading ? 'Updating...' : 'Update Episode'}</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default EditEpisodeForm;
