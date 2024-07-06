import { Link } from 'react-router-dom';

const NavBar = () => {
  return (
    <nav className="navbar" aria-label="main navigation">
      <div className="navbar-brand">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="64"
          height="64"
          viewBox="0 0 100 100"
        >
          <title>Argon</title>
          <g fill="none" stroke="url(#gradient)" strokeWidth="2">
            <circle cx="50" cy="50" r="30" />
            <circle cx="50" cy="50" r="35" />
            <circle cx="50" cy="50" r="40" />
          </g>
          <polygon fill="url(#gradient)" points="50,30 65,70 35,70" />
          <defs>
            <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" style={{stopColor: '#00557F', stopOpacity: '1'}} />
              <stop offset="100%" style={{stopColor: '#00B2A9', stopOpacity: '1'}} />
            </linearGradient>
          </defs>
        </svg>
      </div>

      <a
        href="/settings"
        role="button"
        className="navbar-burger"
        aria-label="menu"
        aria-expanded="false"
        data-target="navbarBasicExample"
      >
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span style={{ display: 'none' }}>Menu</span>
      </a>

      <div id="navbarBasicExample" className="navbar-menu">
        <div className="navbar-start">
          <Link to="/" className="navbar-item">
            Home
          </Link>
          <Link to="/docs" className="navbar-item">
            Documentation
          </Link>
          <Link to="/about" className="navbar-item">
            About
          </Link>
          <div className="navbar-item has-dropdown is-hoverable">
            <Link to="/more" className="navbar-link">
              Admin
            </Link>
            <div className="navbar-dropdown">
              <Link to="/movie/add" className="navbar-item">
                Add movie
              </Link>
              <Link to="/tvshow/add" className="navbar-item">
                Add tvshow
              </Link>
              <hr className="navbar-divider" />
              <Link to="/movie/edit" className="navbar-item">
                Edit movie
              </Link>
              <Link to="/tvshow/edit" className="navbar-item">
                Edit show
              </Link>
              <hr className="navbar-divider" />
              <Link to="/movie/delete" className="navbar-item">
                Delete movie
              </Link>
              <Link to="/tvshow/delete" className="navbar-item">
                Delete tvshow/season/episode
              </Link>
            </div>
          </div>
        </div>

        <div className="navbar-end">
          <div className="navbar-item">
            <div className="buttons">
              <Link to="sign-up" className="button is-primary">
                <strong>Sign up</strong>
              </Link>
              <Link to="log-in" className="button is-light">
                Log in
              </Link>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default NavBar;
