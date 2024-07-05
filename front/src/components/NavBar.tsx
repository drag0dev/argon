import { Link } from 'react-router-dom';

const NavBar = () => {
  return (
    <nav className="navbar" aria-label="main navigation">
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
              <Link to="/series/add" className="navbar-item">
                Add series
              </Link>
              <hr className="navbar-divider" />
              <Link to="/movie/edit" className="navbar-item">
                Edit movie
              </Link>
              <Link to="/series/edit" className="navbar-item">
                Edit show
              </Link>
              <hr className="navbar-divider" />
              <Link to="/movie/delete" className="navbar-item">
                Delete movie
              </Link>
              <Link to="/series/delete" className="navbar-item">
                Delete series/season/episode
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
