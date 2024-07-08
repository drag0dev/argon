import React, { useState } from 'react';
import { signIn } from 'aws-amplify/auth';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const [credentials, setCredentials] = useState({
    username: '',
    password: '',
  });
  const navigate = useNavigate();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCredentials({ ...credentials, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    // Add your API call here to log in the user
    console.log('Login:', credentials);
    try {
        await signIn({username: credentials.username, password: credentials.password, options: {authFlowType: 'USER_PASSWORD_AUTH'}})
        navigate('/')
    } catch (error) {
        alert('error logging in:' +  error)
    }
  };

  return (
    <section className="section">
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2 className="title is-4">Login</h2>

        <div className="field">
          <label className="label" htmlFor="email">Username</label>
          <div className="control">
            <input
              className="input"
              type="username"
              id="username"
              name="username"
              value={credentials.username}
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
