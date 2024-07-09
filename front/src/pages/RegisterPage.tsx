import { signUp } from 'aws-amplify/auth';
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Register = () => {
  const [user, setUser] = useState({
    username: '',
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dateOfBirth: '',
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setUser({ ...user, [name]: value });
  };
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    // Add your API call here to register the user
    console.log('Register:', user);
    try {
        await signUp({
            username: user.username,
            password: user.password,
            options: {
                userAttributes: {
                    email: user.email,
                    'custom:firstName': user.firstName,
                    'custom:lastName': user.lastName,
                    'custom:dateOfBirth': user.dateOfBirth,
                }
            }
        })
        navigate('/')
    } catch (error) {
        alert('error signing up:' +  error)
    }
  };

  return (
    <section className="section">
      <div className="container">
        <form onSubmit={handleSubmit}>
          <h2 className="title is-4">Register</h2>

          <div className="field">
            <label className="label" htmlFor="username">
              Username
            </label>
            <div className="control">
              <input
                className="input"
                type="text"
                id="username"
                name="username"
                value={user.username}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="field">
            <label className="label" htmlFor="email">
              Email
            </label>
            <div className="control">
              <input
                className="input"
                type="email"
                id="email"
                name="email"
                value={user.email}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="field">
            <label className="label" htmlFor="firstName">
              First Name
            </label>
            <div className="control">
              <input
                className="input"
                type="text"
                id="firstName"
                name="firstName"
                value={user.firstName}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="field">
            <label className="label" htmlFor="lastName">
              Last Name
            </label>
            <div className="control">
              <input
                className="input"
                type="text"
                id="lastName"
                name="lastName"
                value={user.lastName}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="field">
            <label className="label" htmlFor="dateOfBirth">
              Date of Birth
            </label>
            <div className="control">
              <input
                className="input"
                type="date"
                id="dateOfBirth"
                name="dateOfBirth"
                value={user.dateOfBirth}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="field">
            <label className="label" htmlFor="password">
              Password
            </label>
            <div className="control">
              <input
                className="input"
                type="password"
                id="password"
                name="password"
                value={user.password}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="field">
            <div className="control">
              <button type="submit" className="button is-primary">
                Register
              </button>
            </div>
          </div>
        </form>
      </div>
    </section>
  );
};

export default Register;
