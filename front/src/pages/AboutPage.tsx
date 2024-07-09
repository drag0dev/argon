import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUser } from '@fortawesome/free-solid-svg-icons';

const AboutPage = () => {
  const teamMembers = [
    { name: "Teodor Đurić", id: "SV67/2021" },
    { name: "Dragoslav Tamindžija", id: "SV47/2021" },
    { name: "Darko Svilar", id: "SV50/2021" }
  ];

  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Our Team</h1>
        <div className="columns is-centered">
          {teamMembers.map((member, index) => (
            <div key={member.id} className="column is-one-third">
              <div className="card">
                <div className="card-content">
                  <div className="media">
                    <div className="media-left">
                      <figure className="image is-48x48">
                        <FontAwesomeIcon icon={faUser} size="3x" className="has-text-primary" />
                      </figure>
                    </div>
                    <div className="media-content">
                      <p className="title is-5">{member.name}</p>
                      <p className="subtitle is-6">{member.id}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default AboutPage;
