import React, { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faEdit } from '@fortawesome/free-solid-svg-icons';
import VideoUpload from './VideoUpload';

const API_URL = process.env.API_URL;

interface VideoMetadata {
  type: string;
  size: number;
  lastModified: Date;
}

interface Movie {
  id: string;
  title: string;
  description: string;
  genres: string[];
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

const EditMovieForm = ({ movie: initialMovie, setEditingMovie }) => {
  const [movie, setMovie] = useState<Movie>(initialMovie);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Initialize form with movie data when component mounts
    setMovie(initialMovie);
  }, [initialMovie]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setMovie({ ...movie, [name]: value });
  };

  const handleArrayInputChange = (e, field) => {
    const { value } = e.target;
    setMovie({
      ...movie,
      [field]: value.split(',').map((item) => item.trim()),
    });
  };

  const handleVideoUpload = (file, metadata: VideoMetaData) => {
    setMovie({
      ...movie,
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
      // PUT request to update the metadata ?
      const response = await fetch(`${API_URL}/movie`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          uuid: movie.id,
          fileType: movie.video.fileType,
          fileSize: movie.video.fileSize,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to update movie');
      }

      const data = await response.json();

      const uploadResponse = await fetch(data.url, {
        method: data.method,
        body: movie.video.file,
        headers: {
          'Content-Type': movie.video.fileType,
        },
      });

      if (!uploadResponse.ok) {
        throw new Error('Failed to upload video file');
      }

      console.log('Video uploaded successfully.');

      // TODO: send separate request for movie details

      alert('Movie updated successfully!');
      setEditingMovie(null); // Close the form upon success
    } catch (error) {
      console.error('Error updating movie:', error);
      alert('Failed to update movie. See the console for more information.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Edit Movie Information</h2>

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
              value={movie.title}
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
              value={movie.description}
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
              value={movie.genres.join(', ')}
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
              value={movie.actors.join(', ')}
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
              value={movie.directors.join(', ')}
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
              <span>{isLoading ? 'Updating...' : 'Update Movie'}</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default EditMovieForm;
