import React, { useState } from 'react';

const Login = () => {
  const [credentials, setCredentials] = useState({
    email: '',
    password: '',
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCredentials({ ...credentials, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    // Add your API call here to log in the user
    console.log('Login:', credentials);
  };

  return (
    <section className="section">
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Login</h2>

        <div className="field">
          <label className="label" htmlFor="email">Email</label>
          <div className="control">
            <input
              className="input"
              type="email"
              id="email"
              name="email"
              value={credentials.email}
              onChange={handleInputChange}
              required
            />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="password">Password</label>
          <div className="control">
            <input
              className="input"
              type="password"
              id="password"
              name="password"
              value={credentials.password}
              onChange={handleInputChange}
              required
            />
          </div>
        </div>

        <div className="field">
          <div className="control">
            <button type="submit" className="button is-primary">
              Login
            </button>
          </div>
        </div>
      </form>
    </div>
    </section>
  );
};

export default Login;
