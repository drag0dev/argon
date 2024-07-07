import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus } from '@fortawesome/free-solid-svg-icons';
import VideoUpload from './VideoUpload';
import { useEffect } from 'react';

const API_URL = process.env.API_URL;

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

  // Function to insert dummy data
  const insertDummyData = (n) => {
    const dummyMovies = [
      {
        id: 1,
        title: 'Inception',
        description:
          'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.',
        genres: ['Action', 'Adventure', 'Sci-Fi'],
        actors: ['Leonardo DiCaprio', 'Joseph Gordon-Levitt', 'Ellen Page'],
        directors: ['Christopher Nolan'],
      },
      {
        id: 2,
        title: 'Interstellar',
        description:
          "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
        genres: ['Adventure', 'Drama', 'Sci-Fi'],
        actors: ['Matthew McConaughey', 'Anne Hathaway', 'Jessica Chastain'],
        directors: ['Christopher Nolan'],
      },
      {
        id: 3,
        title: 'The Dark Knight',
        description:
          'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.',
        genres: ['Action', 'Crime', 'Drama'],
        actors: ['Christian Bale', 'Heath Ledger', 'Aaron Eckhart'],
        directors: ['Christopher Nolan'],
      },
      {
        id: 4,
        title: 'Unexpected Aquatic Humor',
        description:
          'A lighthearted video featuring two different fish, each representing their country in a humorous manner. The first clip showcases a Serbian fish lying on the ground, humorously critiqued for its quality. The second clip features a Bosnian fish "smoking" a cigarette, praised for its swagger.',
        genres: ['Comedy', 'Short', 'Meme'],
        actors: ['Serbian Fish', 'Bosnian Fish'],
        directors: ['Internet Meme Creator'],
      },
    ];

    // Example of using the first dummy movie
    const selectedMovie = dummyMovies[n];
    setMovie({
      ...movie,
      title: selectedMovie.title,
      description: selectedMovie.description,
      genres: selectedMovie.genres,
      actors: selectedMovie.actors,
      directors: selectedMovie.directors,
      video: null,
    });
  };

  // Expose the function to the window object
  useEffect(() => {
    window.insertDummyData = insertDummyData;
    return () => {
      window.insertDummyData = undefined;
    };
  }, []); // Now tracking 'movie' as well, which might be overkill

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
        fileSize: +metadata.size, // this should be illegal
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

      // Step 1: Send movie metadata and get upload URL
      const metadataResponse = await fetch(`${API_URL}/movie`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(movieData),
      });

      if (!metadataResponse.ok) {
        throw new Error('Failed to submit movie metadata');
      }

      console.log('Movie metadata submitted successfully');

      const { url, method } = await metadataResponse.json();

      console.log('Received upload URL:', url);
      console.log('Received upload method:', method);

      // Step 2: Upload the video file directly
      const uploadResponse = await fetch(url, {
        method: method,
        body: movie.video.file,
        headers: {
          'Content-Type': movie.video.fileType,
        },
      });

      if (!uploadResponse.ok) {
        throw new Error('Failed to upload video file');
      }

      console.log('Movie and video uploaded successfully.');

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
