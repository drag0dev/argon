import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUpload } from '@fortawesome/free-solid-svg-icons';

const VideoUpload = ({ onFileUpload, editing }) => {
  const [file, setFile] = useState(null);
  const [metadata, setMetadata] = useState({});

  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
      const fileMetadata = {
        name: selectedFile.name,
        type: selectedFile.type,
        size: selectedFile.size,
        lastModified: new Date(selectedFile.lastModified),
      };
      setMetadata(fileMetadata);
      onFileUpload(selectedFile, fileMetadata); // Call the parent's callback
    }
  };

  return (
    <div className="card pt-4 pb-4">
      <div className="card-content columns is-centered">
        <div className="file column is-narrow">
          <h2 className="title is-4">Upload Video</h2>
          <label className="file-label">
            <input
              className="file-input"
              type="file"
              name="resume"
              onChange={handleFileChange}
            />
            <span className="file-cta">
              <span className="file-icon">
                <FontAwesomeIcon icon={faUpload} />
              </span>
              <span className="file-label"> Choose a file… </span>
            </span>
          </label>
        </div>

        <div
          className="column is-narrow m-2"
          style={{ backgroundColor: 'dimgray', width: '1px', padding: '0' }}
        />

        <div className="column">
          {!file && editing && (
            <h2 className="title is-4">Leave empty to keep the current video</h2>
          )}

          {file && (
            <div>
              <h2 className="title is-4">File Metadata</h2>
              <p>Name: {metadata.name}</p>
              <p>Type: {metadata.type}</p>
              <p>Size: {metadata.size} bytes</p>
              <p>Last Modified: {metadata.lastModified.toLocaleDateString()}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default VideoUpload;
