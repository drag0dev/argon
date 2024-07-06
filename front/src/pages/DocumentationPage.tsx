import React from 'react';
import ReactPlayer from 'react-player';

const DocumentationPage = () => {
  return (
    <section className="section">
      <div className="container">
        Tutorial video:
          <ReactPlayer
            url={'https://www.youtube.com/watch?v=xvFZjo5PgG0'}
            controls={true}
            playing={true}
            className="mb-4"
          />
      </div>
    </section>
  );
};

export default DocumentationPage;
