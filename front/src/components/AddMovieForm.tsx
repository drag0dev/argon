import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus } from '@fortawesome/free-solid-svg-icons';

const AddMovieForm = () => {
  const [movie, setMovie] = useState({
    title: '',
    description: '',
    genres: [],
    actors: [],
    directors: [],
    video: { url: '', format: '' },
  });

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

  const handleVideoInputChange = (e) => {
    const { name, value } = e.target;
    setMovie({ ...movie, video: { ...movie.video, [name]: value } });
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    // TODO: actually send the movie to the server
    console.log('Submitting movie:', movie);
  };

  return (
    <div className="container">
      <h1 className="title is-2">Add New Movie</h1>
      <form onSubmit={handleSubmit}>
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

        <div className="field">
          <label className="label" htmlFor="videoUrl">
            Video URL
          </label>
          <div className="control">
            <input
              className="input"
              type="url"
              id="videoUrl"
              name="url"
              value={movie.video.url}
              onChange={handleVideoInputChange}
            />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="videoFormat">
            Video Format
          </label>
          <div className="control">
            <input
              className="input"
              type="text"
              id="videoFormat"
              name="format"
              value={movie.video.format}
              onChange={handleVideoInputChange}
            />
          </div>
        </div>

        <div className="field">
          <div className="control">
            <button type="submit" className="button is-primary">
              <span className="icon">
                <FontAwesomeIcon icon={faPlus} />
              </span>
              <span>Add Movie</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default AddMovieForm;
