import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus } from '@fortawesome/free-solid-svg-icons';
import VideoUpload from './VideoUpload';

import { API_URL } from '../../config.ts';

const AddMovieForm = () => {
  const [movie, setMovie] = useState({
    title: '',
    description: '',
    genres: [],
    actors: [],
    directors: [],
    video: null,
  });

  const [isLoading, setIsLoading] = useState(false);

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

  const handleVideoUpload = (file, metadata) => {
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
      // Prepare the movie data
      const movieData = {
        title: movie.title,
        description: movie.description,
        genres: movie.genres,
        actors: movie.actors,
        directors: movie.directors,
        video: {
          fileType: movie.video.fileType,
          fileSize: movie.video.fileSize,
          creationTimestamp: movie.video.creationTimestamp,
          lastChangeTimestamp: movie.video.lastChangeTimestamp,
        },
      };

      const url = `${API_URL}/api/movies`;

      // Step 1: Send movie metadata and get upload URL
      const metadataResponse = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(movieData),
      });

      if (!metadataResponse.ok) {
        throw new Error('Failed to submit movie metadata');
      }

      console.log(await metadataResponse.json());

      const { uploadUrl, movieId } = await metadataResponse.json();

      // Step 2: Upload the video file directly
      const uploadResponse = await fetch(uploadUrl, {
        method: 'PUT',
        body: movie.video.file,
        headers: {
          'Content-Type': movie.video.fileType,
        },
      });

      if (!uploadResponse.ok) {
        throw new Error('Failed to upload video file');
      }

      console.log('Movie and video uploaded successfully. Movie ID:', movieId);

      // Clear the form
      setMovie({
        title: '',
        description: '',
        genres: [],
        actors: [],
        directors: [],
        video: null,
      });

      alert('Movie added successfully!');
    } catch (error) {
      console.error('Error adding movie:', error);
      alert('Failed to add movie. See the console for more information.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Movie Information</h2>

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

        <VideoUpload onFileUpload={handleVideoUpload} />

        <div className="field">
          <div className="control">
            <button
              type="submit"
              className="button is-primary"
              disabled={isLoading || !movie.video}
            >
              <span className="icon">
                <FontAwesomeIcon icon={faPlus} />
              </span>
              <span>{isLoading ? 'Uploading...' : 'Add Movie'}</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default AddMovieForm;
